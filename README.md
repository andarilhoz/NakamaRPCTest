# README - Nakama Test

## Solution Explanation

### How to Run?

To run the application, simply execute the following command from the root folder:

```shell
docker-compose up -d
```

## Organization

The code structure went through different approaches during development. Initially, files were separated by context in each package. However, i realized that maybe that would be an overengineering, so i decided to keep all files in the same folder to improve understandability and deployment since it is a simple application. Functions were kept small, following the principle of single responsibility, which resulted in multiple `.go` files in exchange of a cleaner code. Additionally, a pointer was included inside the payload request and payload response to fulfill the requirement of "null" for the hash and content as mentioned in the test.

## Database Storage

A new table called "requests" was created to store individual successful requests. This decision was made to track cases where the hashes didn't match and to provide an audit trail of the requested operations. As the content is already stored as a file in the server, i didn't include it in the table. 

## Thoughts and Ideas about the Task

- The specific use case for this function in the real world is not clear yet for me. It could potentially be used for retrieving saves or managing configurations, but i'm not sure.
- The logic of the task was not overly challenging, but it provided an opportunity to explore Nakama and experience coding in Go, wich was really cool.
- Mocking Nakama functions (such as NakamaModule, Logger, DB, etc.) posed difficulties, so interfaces were introduced early in the code to decouple the function dependencies. As a result, the main function and the initial part of the RPC were not covered by tests.
- I considered changing the error return format to JSON and adopting a message-based approach similar to REST. However, due to Nakama documentation recommending the use of custom error messages, the implementation followed their guidelines.

## Possible Improvements with More Time

Given additional time, the following improvements could be made:

- Perform a check on the size of the JSON file before reading it, potentially slicing the reading process by lines to handle massive file sizes.
- Implement integration tests, particularly focusing on the database functionality.
- Create a dedicated error handling function, optionally allowing the sending of JSON instead of Nakama error messages.

## Test Feedback

I felt that the test could have been a bit more complex and focused on real-life scenarios, such as implementing a battle pass or a reward chest system. 
It also would have been helpful to have a use case for the function in the Story section.

Nevertheless, I greatly appreciated the opportunity to participate in this test and learn a little more about pre-made backend servers and Go itself.
