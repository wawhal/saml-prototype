#SAML 2.0 prototype

For testing, follow the given steps (little tedious I know):

1. Build the deps and run the server on port 8080. `go build main.go && ./main.go` from `server` directory. Or whatever's your favourite way to run go files.
2. Run the client on port 3000. `npm start` from `client` directory.
3. Point `localhost:8080` to `server.saml.test` and `localhost:3000` to `ui.saml.test`. This is important for sharing cookies between UI and Server. I did it the following way; if you know a better way, feel free to use that.

	#. Add the following site to your nginx conf. Create a file called `mysite` in `/etc/nginx/sites-available/` and add the following content:

	```
	server {
	    listen 80;
	    server_name ui.saml.test;

	    location / {
	        proxy_pass http://localhost:3000;
	        proxy_set_header Host $host;
	        proxy_set_header X-Real-IP $remote_addr;
	        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	        proxy_set_header X-Forwarded-Proto $scheme;
	    }
	}

	server {
	    listen 80;
	    server_name server.saml.test;

	    location / {
	        proxy_pass http://localhost:8080;
	        proxy_set_header Host $host;
	        proxy_set_header X-Real-IP $remote_addr;
	        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	        proxy_set_header X-Forwarded-Proto $scheme;
	    }
	}
	```

	#. Resolve these new domains in `/etc/hosts`. Add these two lines after the localhost resolution:

	```language
	127.0.0.1	ui.saml.test
	127.0.0.1	server.saml.test
	```

	#. Enable this newly added site:

	```language
	ln -s /etc/nginx/sites-available/mysite /etc/nginx/sites-enabled/mysite 
	```

	#. Test config and restart nginx

	```
	sudo service nginx configtest
	sudo service nginx restart
	```


4. Go to `ui.saml.test` and login using `user@saml.test` and `Hasura123@`
