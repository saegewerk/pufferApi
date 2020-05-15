	for i in {1..5}
	do
	 newman run /Users/dgstoehl/Documents/Postman/pufferApi.postman_collection.json -n 10000  &
	done
	#--reporter-cli-silent