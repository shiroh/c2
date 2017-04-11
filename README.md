# Carou

### API
The package serve the http request, including `createTopic`, `getTopics` and `voteUpdate`
#### createTopic
Create a new topic and write the the storage. It will initialize an UUID as a ID.
#### getTopics
Return at most page size topics sorted by votes
#### updateVote
Increase/Decrease the votes according to the topic id.


### Cache
It is a layer to cache the topics in the memory.
Each Update op will refresh the value in the cache as well as the storage.
Each Get op will get all topics from cache OR from the load from storage if they are expired. `TTl` SHOULD be configurable, but it is hard coded as 1 second.
 
### Errors
Define all the errors share across the project.

### Modles
Define the data structure to describe the data used in the project.

### Storage
It is a layer to put the topics to the storage and load from it.
In the demo it is a concurrent hashmap for demo purpose.
But I'd like to use a queue in future version to update the vote changes. It would aggregate vote changes by batch to reduce the dp ops.

### Utils
Define the data structure, functions and modules used in the project

### Vendors
Contains the dependencies

### Views
Contains the html pages.
