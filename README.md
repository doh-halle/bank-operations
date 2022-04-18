This is the source code for a bank app, with API endpoints for performing the basic bank transactions listed below.

    Create a new customer account
    Deposit money into a specific account
    Withdraw money from a specific account
    Check the balance of a specific account

The app was written in Golang and it enterracts with a postgresql database at the backend. Regarding my level of experience, I have been working with Go for two years. The transactions were tested using the web API testing tool - Postman

The installation steps assume that you have docker and postman installed on your test computer. run the following commands to download, build and run the two containers necessary for this test, One is a postgres dababase called bo-db and the other is a golang api app.

$ docker build -t hallecraft/bo-db:version1 $ docker build -t hallecraft/bo-api:version1

$ docker run -p 5400:5432 hallecraft/bo-db $ docker run -p 7000:7000 hallecraft/bo-api

To ensure the containers are up and running, run the command below, you should get a list with both containers and status: up $ docker ps

Check the available docker networks using the command below, you should find bo-network in the list $ docker network ls

Connect the API container - hallecraft/bo-api into the bo-network after copying its container id from the results of the docker ps command, with the command below to connect. $ docker network connect bo-network container_id

The system is now ready for you to begin API tests. Launch your postman app.

Test 1 - Account creation test.

    place the link - localhost:7000/api/create_customer - in address or URL bar
    Choose POST request, from frequest methods next to the address bar
    Choose raw and json within the body of the request
    Paste in the following json formatted data in the body of the request { "firstname": "Barack", "lastname": "Obama", "email": "bobama@usa.gov", "phonenumber": "5124684442", "occupation": "President", "customercity": "W-DC" }
    Click the Blue send button

Expected response should be

    200 for response code
    response body - { "id": 5, "message1": "Account created successfully, Thank you for using Mondu!" }

Test 2 - Deposit Cash into specific account test.

    place the link - localhost:7000/api/deposit_cash/1 - in address or URL bar
    Choose POST request, from frequest methods next to the address bar
    Choose raw and json within the body of the request
    Paste in the following json formatted data in the body of the request

{ "accountnumber": "1", "mediumoftransaction": "Transfer", "transactionamount": "28000" }

    Click the Blue send button

Expected response should be

    200 for response code
    Responds body - { "id": 1, "message1": "New deposit completed successfully, Amount - 28000", "message2": " - Your new account balance is - 28000" }

Test 3 - Withdraw Cash from specific account test.

    place the link - localhost:7000/api/withdraw_cash/1 - in address or URL bar
    Choose POST request, from frequest methods next to the address bar
    Choose raw and json within the body of the request
    Paste in the following json formatted data in the body of the request { "accountnumber": "1", "mediumoftransaction": "Cash", "transactionamount": "10000" }
    Click the Blue send button

Expected response should be

    200 for response code
    Responds body { "id": 1, "message1": "New withdrawal completed successfully, Amount - 10000", "message2": " - Your new account balance is - 18000" }
    If transaction amount was equal or more than the current balance, response should be "Insufficient Funds!"

Test 4 - Check Account Balance for specific account id test.

    place the link - localhost:7000/api/account_balance/1 - in address or URL bar
    Choose GET request, from frequest methods next to the address bar
    Click the Blue send button

Expected response should be

    200 for response code
    Response body { "accountnumber": "1", "customerid": "1", "accountbalance": "43000", "balancedate": "2022-04-18T00:00:00Z" }

If ID spedified doesn't exist in the database, response body display "Account Non-existent!"
