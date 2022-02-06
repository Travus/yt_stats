# YouTube Stats
Welcome to the YouTube Status repository. YouTube Status is a REST API designed to be used as a middle man to more easily get relevant information and statistics on different parts of YouTube.


**Current Features:**
* Get statistics and information relating to up to 50 channels at once.
* Get statistics and information relating to up to 50 playlists at once.
    * Also get information and statistics on the contained videos with the same request.
    * Automatically let the REST API calculate total statistics, averages, and more.
* Get statistics and information relating to up to 50 videos at once.
    * The REST API can calculate total statistics for these unrelated videos as well.
* Get information on up to 50 livestreams at once.
* Get all comment and replies on a video in 1 request, no more pagination and fishing for replies.
    * These comments and replies can also be extensively filtered by author and message content.
    * Supports multiple additive and reductive filters, toggleable case sensitivity, and more.
* Get live chat messages and events from an active livestream.
    * This also shows other events such as SuperChats and Memberships.
* YouTube Stats lets you track your quota usage by telling you it's usage.
* A status endpoint to see if the REST API and YouTube API is operational.

Once set up you can use it with all your other apps. By letting a serialized REST API handle these things for you, you no longer have to implement the same functionality in all of your apps that need similar things, and adding new functionality to the REST API makes it available for all your apps with minimal effort.

This REST API takes YouTube API keys as a parameter in the request header, which means that multiple people or applications can use it with separate API keys without sharing your YouTube quota. It's designed to be as light on the quota usage as possible, and reports back how much of your quota it used so you can more easily track it.

## Setup
This repository comes with a Dockerfile in order to run it using Docker. If you do not have Docker already installed then you can follow [Docker's installation guide](https://docs.docker.com/docker-for-windows/install/) in order to install it. There is also a docker-compose file included for even simpler deployment.  

### Deploy with docker-compose
Start by cloning this repository, and navigating into it from your command line.

If you wish to use TLS, then edit the `docker-compose.yml` file, and uncomment the commented-out lines by removing the `#` symbols from the beginning of them. Then replace `example.com` with the address YouTube Stats will be reachable under. If you do not wish to use TLS, you can skip this step.

Run `docker-compose up -d` and docker-compose will start YouTube Stats for you.

### Deploy with Docker
Start by cloning this repository, and navigating into it from your command line.
 
Run `docker build -t yt_stats:latest .` to create a docker image of YouTube Stats
> **Note:** yt_stats:latest can be replaced with any other name you'd like to give the image.

Run one of the following commands, based on if you wish to use TLS or not. If you chose a different name for the image in the previous command, you will have to replace it here too.

With TLS:  
`docker run -d -p 80:8080 -p 443:8081 -v cert-cache:/app/cert-cache -e "tls_address=..." yt_stats:v1`
> **Note:** Replace `...` with the address YouTube Stats will be reachable under.

Without TLS:  
`docker run -d -p 80:8080 yt_stats:v1`

If both commands worked as they should, you'll have a running instance of YouTube Stats now. You can test this by opening `YOUR_ADDRESS/ytstats/v1/` in your browser, and you should see some text indicating that you have reached the YouTube Stats REST API.

All you need to do now is to [get your YouTube API key](https://github.com/Travus/yt_stats/wiki#getting-a-youtube-api-key) and read up on what the different endpoints return. This is listed in the [wiki](https://github.com/Travus/yt_stats/wiki) attached to this repository.

## Contact
If you have any questions, needs, or requests, feel free to contact me!  
I'm mostly active on Discord, but you can reach me on Twitter too.  
Alternatively, you can also create an issue on this repository, and I will get back to you soon.

Discord: Travus#8888  
Twitter: [@RealTravus](https://twitter.com/RealTravus)
