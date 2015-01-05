dumbflake
=========

Simple UDP daemon to help with issuing monotonically increasing UID's for 
newly registered users.

Workflow
--------
1. Server receives UDP packet with login inside.
2. Server checks if there is a reserved ID for specified login.
2a. If there is a reservation, assign reserved ID and do not modify global counter.
2b. If there is no reservation found, increase global counter and assign increased value.
3. Server replies with UDP packet, structured as follows: "<ID>:<login>\n"
	
Usage
-----

    dumbflake -port=19229 -bind=127.0.0.1
	
