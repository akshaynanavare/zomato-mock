# shortest-time-of-deliveries

## How to run this program?

1. Install `go 1.22` version
2. Run `make run` command
3. Hit this curl URL to get the shortest-time-of-deliveries of default values
    `curl --location --request GET 'http://127.0.0.1:8080/delivery/shortest-time/aman'`
4. If you want to change the default values, check out router/defaults.go file and change order's data accordingly.


### How alogithm works ?
Step 1:
    Build graph data structure from availble restaurants and driver's current location
Step 2:
    Set maximum weight to an edge between driver and restaurant's node => max(time taken to reach restaurant, avg preparation time taken by restaurant)
Step 3:
    Build heap from available nodes
Step 4:
    Choose the minimum weight from the available nodes
Step 5:
    If the given node is restaurant then update graph by adding the restaurant's customer because we can not visit customers without picking meal from restaurant
Step 6:
    Add customer node's edge in all unvisited nodes because the next node might be some other restaurant and in that case our first customer wont be reached
Step 7:
    Get all the neighbors of the restaurant and compute the heap again
Step 8:
    Update heap to the processing heap
Step 9:
    Repeat the process from step 4 until our heap is empty
Step 10:
    return last heap node with having time and path.


### Optimization scope
1. Above algorithm chooses the greedy approach where it might not give the most optimal solution.
2. Using the DP by visiting all possible paths will be the optimized solution.
3. Improvment scope in memory utilization.

