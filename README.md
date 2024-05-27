# Deliver orders in the possible shortest time.

This project calculates the shortest time of deliveries based on the given data of restaurants, customers and driver.

## Prerequisites

- Install Go version 1.22

## Installation

1. Ensure you have Go 1.22 installed on your machine. You can download and install it from [the official Go website](https://golang.org/dl/).

2. Clone the repository and navigate to the project directory.

3. Run the following command to install dependencies and start the application:
    ```sh
    make run
    ```

## Usage

To get the shortest time of deliveries with default values, use the following CURL command:
```sh
curl --location --request GET 'http://127.0.0.1:8080/delivery/shortest-time/aman'
```
If you want to change the default values, check out the router/defaults.go file and update the order's data accordingly.

## How the Algorithm Works
1. **Build Graph Data Structure**: Create a graph from available restaurants and the driver's current location.
2. **Set Maximum Weight**: Assign the maximum weight to an edge between the driver and a restaurant node. This weight is the maximum of the time taken to reach the restaurant and the average preparation time taken by the restaurant.
3. **Build Heap**: Create a heap from available nodes.
4. **Choose Minimum Weight Node**: Select the minimum weight node from the heap.
5. **Update Graph**: If the selected node is a restaurant, update the graph by adding the restaurant's customer. This is because a customer cannot be visited without picking up the meal from the restaurant first.
6. **Add Customer Node's Edge**: Add edges for the customer node in all unvisited nodes because the next node might be another restaurant, and in that case, the first customer won't be reached.
7. **Get Neighbors**: Retrieve all neighbors of the restaurant and compute the heap again.
8. **Update Heap**: Update the processing heap.
9. **Repeat Process**: Repeat the process from step 4 until the heap is empty.
10. **Return Result**: Return the last heap node containing the time and path.

## Time and Space Complexity
1. Time = `O(V + E * log(V))` V : total nodes, E : total edges
2. Space = `O(V + E)` 

## Oprimzation scope
1. The above algorithm uses a greedy approach and might not provide the most optimal solution.
2. Using Dynamic Programming (DP) by visiting all possible paths will yield a more optimized solution.
3. There is scope for improvement in memory utilization.
4. Unit test cases should be considered due to time constraints
5. Handle multiple customers from same restaurant. (current assumption 1:1 relation between customer and restaurant)