# DDOS automation tool for Google Cloud Platform

### Description
Glory for Ukraine🇺🇦


### How to start
1. Register account on GCP (you will get $300 for free) https://cloud.google.com/free
2. Create new project https://cloud.google.com/resource-manager/docs/creating-managing-projects or use existing one
3. Create and download a new Api Key https://cloud.google.com/docs/authentication/api-keys or use existing one
4. Copy your project id from Project Info
5. Download one of binaries from Releases for your computer architecture
6. Just use


### How to use
1. Help ```./gcp-ddos --help```
1. Create instances: 
```
    ./gcp-ddos --command=create --pid=<your_project_id> \
    --key=<path_to_api_key.json> \
    --url=<target_url> \
    -d=<duration>
```
2. Delete instances:
```
    ./gcp-ddos --command=delete --pid=<your_project_id> \
    --key=<path_to_api_key.json>
```
3. Stop instances:
```
    ./gcp-ddos --command=stop --pid=<your_project_id> \
    --key=<path_to_api_key.json>
```


### Upcomming
1. Number of instances selection
2. Zones selection
3. Startup script customization
4. List all instances

