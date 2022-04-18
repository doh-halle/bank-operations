CREATE TABLE customer
   (
       customerid BIGSERIAL NOT NULL,
       firstname VARCHAR(20) NOT NULL,
       lastname VARCHAR(20) NOT NULL,
       email VARCHAR(30)NOT NULL,
       phonenumber VARCHAR(20)NOT NULL,
       occupation VARCHAR(20),
       customercity VARCHAR(20) not NULL,
       createdat DATE DEFAULT current_timestamp,
      CONSTRAINT customer_customerid_pk PRIMARY KEY(customerid)  
   );
 
  CREATE TABLE accountbalance (
      accountnumber BIGSERIAL NOT NULL,
      customerid BIGSERIAL NOT NULL,
      accountbalance NUMERIC,
      balancedate DATE DEFAULT current_timestamp,
      CONSTRAINT accountbalance_accountnumber_pk PRIMARY KEY(accountnumber),
      CONSTRAINT accountbalance_customerid_fk FOREIGN KEY(customerid) REFERENCES customer(customerid)
   );

  CREATE TABLE transactionhistory (
   transactionnumber BIGSERIAL not null,
      accountnumber BIGSERIAL NOT NULL,
      transactiondate DATE not null DEFAULT current_timestamp,
      currency VARCHAR(10) NOT NULL,
      mediumoftransaction VARCHAR(20) not NULL,
      transactiontype VARCHAR(20) not null,
      transactionamount NUMERIC,
      CONSTRAINT transactionhistory_transactionnumber_pk PRIMARY KEY(transactionnumber),
      CONSTRAINT transactionhistory_accountnumber_fk FOREIGN KEY(accountnumber) REFERENCES accountbalance(accountnumber)
   );