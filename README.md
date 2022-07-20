# Web Benchmark

Most web benchmarks are deceiving developers. Is it "real" to benchmark a language and a framework without considering their real world usage? When did you last create an application which did not rely on some kind of database or I/O to achieve its goal? Applications often use databases and other types of I/O (disk, internet, etc.) to deliver the expected result. If it's performance (a.k.a. requests per minute and result latency) that you're looking for, you __cannot__ disconsider the influence that a database will have on your benchmark.

For examples, `fastify` claims to be much faster framework than `express` according to their [benchmark](https://www.fastify.io/benchmarks/). But does it hold this performance when both frameworks are tested with databases? Now, imagine that developer squad, based on the benchmark provided by `fastify`, decides to refactor their applications. How would they look like if in the end, after weeks or months of work, both systems deliver the same performance?

To be very clear, I'm not disregarding the effort make to create lightweight and fast languages, runtime and frameworks. These benchmarks provided by them are the consequence of all their effort to fine tune their code in order to achieve better results. All I'm saying is that these "raw benchmarks" __alone__ do not provide real world usage for software engineers who usually have to rely on databases, HTTP requests, etc. to deliver value to their customer. Taking action without this in mind will likely lead you and your team to take wrong decisions.

## Methodology

All the implementations must respect the fallowing rules:

1. The server __must__ listen on port 3000;
2. The benchmark endpoint __must__ be the `GET /`;
3. It should have 100 connections pool size;
4. Benchmarks without databases __must__ return `{"status": "ok"}`;
5. Benchmarks with databases __must__ execute the fallowing commands
    a. Create and persist the user with name `Test!` and age `1`;
    b. Load the user from the database with the `id` returned by the previous command;
    c. Update the user's age to 2 with (e.g. `user.age += 1` or `user.age = user.age + 1`) and save it;
    d. Delete the user entry using its `id`;
    e. The endpoint must return a `200` with the JSON representation of the user entry (e.g. `{"id": 1, "name": "Test!", "age": 2}`).

## Running Locally

### Setting up Kubernetes

- Start minikube and set usage of docker images

```shell
minikube start
eval $(minikube docker-env)
```

- Setup `kubectl`

```shell
minikube -p minikube docker-env
```

- Setup `.env` CLUSTER_IP variable with the result of `minikube ip`

- Enable `metrics-server` (Not yet being used)

```
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type 'json' -p '[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```

### Running the benchmark

```shell
npm i
npm start
```

## Debugging

Useful scripts to debug the benchmark.

- Access the database directly

```shell
kubelet exec -it DatabasePodName -- psql -U user -d mydb
```

- Forwards a pod port in order to analise its operation.

```shell
kubectl port-forward PodName 3000:3000
```
