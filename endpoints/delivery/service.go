package delivery

import (
	"time"

	constants "github.com/akshaynanavare/zomato-mock/constants"
	graph "github.com/akshaynanavare/zomato-mock/graph"
	model "github.com/akshaynanavare/zomato-mock/models"
	repository "github.com/akshaynanavare/zomato-mock/repository"
	utils "github.com/akshaynanavare/zomato-mock/utils"

	"log"
)

type Service interface {
	GetShortestPathOfActiveOrders(deliveryPartner string) (time.Time, []string, error)
}

type service struct {
	Order           repository.Order
	DeliveryPartner repository.DeliveryPartner
}

func NewService(
	order repository.Order, deliveryPartner repository.DeliveryPartner,
) Service {
	return &service{
		Order:           order,
		DeliveryPartner: deliveryPartner,
	}
}

// GetShortestPathOfActiveOrders : Handler for get shortest path of active orders of delivery partner
func (s *service) GetShortestPathOfActiveOrders(deliveryPartner string) (time.Time, []string, error) {
	// get orders of the delivery partner
	orders, err := s.Order.GetOrdersByID(deliveryPartner)
	if err != nil {
		log.Printf("failed to get orders from DB err: %v", err)
		return time.Time{}, nil, err
	}

	// get delivery partner object
	partner, err := s.DeliveryPartner.GetDeliveryPartnerByID(deliveryPartner)
	if err != nil {
		log.Printf("failed to get delivery partner from DB err: %v", err)
		return time.Time{}, nil, err
	}

	// find the shortest time and path
	ts, path := findShortestPath(orders, partner)

	// convert to duration
	duration := time.Duration(ts * float64(time.Minute))

	// return the time and path
	return time.Now().Add(duration), path, nil
}

// findShortestPath : it returns the shortest path time of delivery partner.
func findShortestPath(orders []*model.Order, partner *model.DeliveryPartner) (float64, []string) {
	var (
		travelledTime         float64
		visited               = map[string]bool{}
		restaurantCustomerMap = map[string][]*model.Customer{}
		restaurantMap         = make(map[string]*model.Restaurant)
		top                   utils.Path

		// build graph
		orderGraph, nodes = NewGraph(partner, orders)
	)

	for _, v := range orders {
		if _, ok := restaurantCustomerMap[v.Restaurant.ID]; !ok {
			restaurantCustomerMap[v.Restaurant.ID] = []*model.Customer{}
		}

		restaurantCustomerMap[v.Restaurant.ID] = append(restaurantCustomerMap[v.Restaurant.ID], v.Customer)
		restaurantMap[v.Restaurant.ID] = v.Restaurant
	}

	visited[partner.ID] = true

	// build heap for restaurants
	h := utils.NewHeap()

	// add each restaurant into the heap with weight from max(time to reach restaurant or restaurant prep time)
	for _, n := range orderGraph.GetEdges(partner.ID) {
		maxTime := (n.Weight / constants.AvgBikeSpeed) * constants.Int60
		if maxTime < float64(restaurantMap[n.Node.ID].AvgPrepTime) {
			log.Println("restaurant time is greater than travel time : ", maxTime, "node", n.Node.ID)
			maxTime = float64(restaurantMap[n.Node.ID].AvgPrepTime)
		}

		h.Push(utils.Path{Value: maxTime, Nodes: []string{partner.ID, n.Node.ID}})
	}

	for h.Len() > 0 {
		top = h.Pop()
		heap := utils.NewHeap()
		currNode := top.Nodes[len(top.Nodes)-1]

		travelledTime = top.Value

		log.Printf("[DEBUG] travelled time so far : %f, currNode : %s\n", travelledTime, currNode)

		// add customer in graph if curr node is restaurant
		if customers, ok := restaurantCustomerMap[currNode]; ok {
			for _, c := range customers {
				customerNode := &graph.Node{
					ID:       c.ID,
					Location: c.Location,
				}
				orderGraph.AddEdgeToUnvistitedNodes(nodes, visited, customerNode)
			}
		}

		// prepare next neighbours to visit
		for _, n := range orderGraph.GetEdges(currNode) {
			// if visited then skip
			if visited[n.Node.ID] {
				continue
			}

			// get maxTime by so far by travelling to the node + consumed time
			maxTime := travelledTime + ((n.Weight / constants.AvgBikeSpeed) * constants.Int60)

			// if maxTime is greater than restaurant prep time then use restaurant prep time
			if val, ok := restaurantMap[n.Node.ID]; ok {
				if maxTime < float64(val.AvgPrepTime) {
					maxTime = float64(val.AvgPrepTime)
				}
			}

			heap.Push(utils.Path{Value: maxTime, Nodes: append(top.Nodes, []string{n.Node.ID}...)})
		}

		visited[currNode] = true
		h = heap
	}

	return top.Value, top.Nodes
}

func NewGraph(p *model.DeliveryPartner, orders []*model.Order) (*graph.Graph, map[string]*graph.Node) {
	deliveryPartner := &graph.Node{
		ID:       p.ID,
		Location: p.CurrentLocation,
	}

	nodesMap := map[string]*graph.Node{}
	nodesMap[p.ID] = deliveryPartner

	orderGraph := graph.Graph{
		AdjacencyList: make(map[string][]graph.Edge),
	}

	var prev *graph.Node

	for _, o := range orders {
		currNode := graph.Node{
			ID:       o.Restaurant.ID,
			Location: o.Restaurant.Location,
		}

		if prev != nil {
			orderGraph.AddEdge(prev, &currNode, utils.CalculateDistance(prev.Location, o.Restaurant.Location))
		}

		orderGraph.AddEdge(deliveryPartner, &currNode, utils.CalculateDistance(p.CurrentLocation, o.Restaurant.Location))
		prev = &currNode
		nodesMap[currNode.ID] = &currNode
	}

	return &orderGraph, nodesMap
}
