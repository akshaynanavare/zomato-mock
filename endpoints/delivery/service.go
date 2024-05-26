package delivery

import (
	"fmt"
	"time"

	constants "github.com/akshaynanavare/shortest-time/constants"
	graph "github.com/akshaynanavare/shortest-time/graph"
	model "github.com/akshaynanavare/shortest-time/models"
	repository "github.com/akshaynanavare/shortest-time/repository"
	utils "github.com/akshaynanavare/shortest-time/utils"

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

func (s *service) GetShortestPathOfActiveOrders(deliveryPartner string) (time.Time, []string, error) {
	orders, err := s.Order.GetOrdersByID(deliveryPartner)
	if err != nil {
		log.Printf("failed to get orders from DB err: %v", err)
		return time.Time{}, nil, err
	}

	partner, err := s.DeliveryPartner.GetDeliveryPartnerByID(deliveryPartner)
	if err != nil {
		log.Printf("failed to get delivery partner from DB err: %v", err)
		return time.Time{}, nil, err
	}

	ts, path := findShortestPath(orders, partner)
	// log.Println("path : ", path, "time : ", (distance/constants.AvgBikeSpeed)*60, "min")
	duration := time.Duration(ts * float64(time.Minute))

	return time.Now().Add(duration), path, nil
}

func findShortestPath(orders []*model.Order, partner *model.DeliveryPartner) (float64, []string) {
	orderGraph, nodes := NewGraph(partner, orders)
	visited := map[string]bool{}
	restaurantCustomerMap := map[string][]*model.Customer{}
	restaurantMap := make(map[string]*model.Restaurant)

	var top utils.Path

	for _, v := range orders {
		if _, ok := restaurantCustomerMap[v.Restaurant.ID]; !ok {
			restaurantCustomerMap[v.Restaurant.ID] = []*model.Customer{}
		}

		restaurantCustomerMap[v.Restaurant.ID] = append(restaurantCustomerMap[v.Restaurant.ID], v.Customer)
		restaurantMap[v.Restaurant.ID] = v.Restaurant
	}

	// fmt.Println("original graph : ")
	// orderGraph.PrintGraph()
	visited[partner.ID] = true
	h := utils.NewHeap()
	for _, n := range orderGraph.GetEdges(partner.ID) {
		maxTime := (n.Weight / constants.AvgBikeSpeed) * constants.Int60
		if maxTime < float64(restaurantMap[n.Node.ID].AvgPrepTime) {
			log.Println("restaurant time is greater than travel time : ", maxTime, "node", n.Node.ID)
			maxTime = float64(restaurantMap[n.Node.ID].AvgPrepTime)
		}

		h.Push(utils.Path{Value: maxTime, Nodes: []string{partner.ID, n.Node.ID}})
	}

	var travelledTime float64

	for h.Len() > 0 {
		top = h.Pop()
		heap := utils.NewHeap()
		currNode := top.Nodes[len(top.Nodes)-1]

		travelledTime = top.Value

		log.Printf("[DEBUG] travelled time so far : %f, currNode : %s\n", travelledTime, currNode)

		// add customer in graph if curr node is restaurant
		if customers, ok := restaurantCustomerMap[currNode]; ok {
			for _, c := range customers {
				fmt.Println("adding customer in graph : ", c)
				customerNode := &graph.Node{
					ID:       c.ID,
					Location: c.Location,
				}
				orderGraph.AddEdgeToUnvistitedNodes(nodes, visited, customerNode)
			}

			log.Println("after adding customer in graph : ", customers)
			orderGraph.PrintGraph()
		}

		for _, n := range orderGraph.GetEdges(currNode) {
			if visited[n.Node.ID] {
				continue
			}
			maxTime := travelledTime + ((n.Weight / constants.AvgBikeSpeed) * constants.Int60)

			if val, ok := restaurantMap[n.Node.ID]; ok {
				if maxTime < float64(val.AvgPrepTime) {
					maxTime = float64(val.AvgPrepTime)
				}
			}

			fmt.Println("Got edge for source : ", currNode, " dest : ", n.Node.ID, "maxTime : ", int64(maxTime), "km : ", n.Weight)

			heap.Push(utils.Path{Value: maxTime, Nodes: append(top.Nodes, []string{n.Node.ID}...)})
		}

		visited[currNode] = true
		h = heap
	}

	fmt.Println("final graph : ")
	// orderGraph.PrintGraph()

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

		// customerNode := graph.Node{
		// 	ID:       o.Customer.ID,
		// 	Location: o.Customer.Location,
		// }

		if prev != nil {
			orderGraph.AddEdge(prev, &currNode, utils.CalculateDistance(prev.Location, o.Restaurant.Location))
		}

		orderGraph.AddEdge(deliveryPartner, &currNode, utils.CalculateDistance(p.CurrentLocation, o.Restaurant.Location))
		// graph.AddEdge(&currNode, &customerNode, utils.CalculateDistance(currNode.Location, customerNode.Location))
		prev = &currNode

		nodesMap[currNode.ID] = &currNode
		// nodesMap[customerNode.ID] = &customerNode
	}

	return &orderGraph, nodesMap
}