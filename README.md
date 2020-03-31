# API to Proxy/Cache requests to ML API

This simple API receives requests to /items/$ITEM_ID route and tries to resolve the Item locally and if not possible performs a remote request to ML API.

<h3>BUILD / RUN</h3>


Build docker image
```cmd
docker build -t <img_name> .
```

Run
```cmd
 docker run -p <port>:8080 -e DB_CONN_STR='<db_user>:<db_pass>@tcp(<db_host>:<db_port>)/<db_name>' <img_name>
```

<h3>Important!</h3>

If no ***DB_CONN_STR*** is provided, the system resolves all requests remotly, and ***/health*** route resturns 503 (Service Unavailable) error.

<p>When proxy service is running the following routes are available.</p>
<ul>
  <li>/items/$ITEM_ID</li>
  <li>/health</li>
</ul>

<h2>Further improvements</h3>

At this moment this service is resilient to DB server fault, in that situation it proxies all requests directly to ML remote API. To prevent DB server faults we can build a Master / Slave DB server architecture.

***/health*** route returns a resume of the system behaviour grouped by minute based on requests logged to DB. In a scenario with many requests per second, this process could be very heavy. To make an improvement it is possible to create snapshots requests using a scheduled job, it is noted that we can lose all the snapshots and rebuild them using the requests table wich is our unique source of truth. Instead of retrieving large volumes of data is also reasonable to add parameters to ***/health*** route such as ***since*** parameter, to fetch only a subset of data.

***/items/$ITEM_ID*** In order to manage large number of items and more complex queries, it is possible to create shardings based for example in the location of an Item (address, geolocation ...). To achive this we can use a NoSQL database to create a hastable to map ML ID to internal ID.
