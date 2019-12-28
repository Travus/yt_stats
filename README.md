# YouTube Stats
Welcome to the YouTube Status repository. YouTube Status is a REST API designed to be used as a middle man to more easily get relevant information and statistics on different parts of YouTube.


**Current Features:**
* Get statistics and information relating to up to 50 channels at once.
* Get statistics and information relating to up to 50 playlists at once.
    * Also get information and statistics on the contained videos with the same request.
    * Automatically let the REST API calculate total statistics, averages, and more.
* Get statistics and information relating to up to 50 videos at once.
    * The REST API can calculate total statistics for these unrelated videos as well.
* Get all comment and replies on a video in 1 request, no more pagination and fishing for replies.
    * These comments and replies can also be extensively filtered by author and message content.
    * Supports multiple additive and reductive filters, toggleable case sensitivity, and more.
* YouTube Stats lets you track your quota usage by telling you it's usage.
* A status endpoint to see if the REST API and YouTube API is operational.

Once set up you can use it with all your other apps. By letting a serialized REST API handle these things for you, you no longer have to implement the same functionality in all of your apps that need similar things, and adding new functionality to the REST API makes it available for all your apps with minimal effort.

This REST API takes YouTube API keys as a parameter, which means that multiple people or applications can use it with separate API keys without sharing your YouTube quota. It's designed to be as light on the quota usage as possible, and reports back how much of your quota it used so you can more easily track it.

## Setup
This repository comes with a Dockerfile in order to run it using Docker. If you do not have Docker already installed then you can follow [Docker's installation guide](https://docs.docker.com/docker-for-windows/install/) in order to install it.  
Once you have Docker installed you can simply clone this repository and build a Docker image by navigating into the repository and running the following command:  
`docker build -t yt_stats:v1`
> Note: yt_stats:v1 can be replaced with any other name you'd like to give the image.

Once the docker image is created, you can start a docker instance of it with the following command:  
`docker run -d -p 80:8080 yt_stats:v1`
> Note: This will bind port 80 to redirect to the REST API on port 8080, which is where the REST API is listening. If you have other things running on port 80, things might break. Additionally, if you chose a different name for the image in the previous command, you will have to replace it here too.

If both commands worked as they should you'll have a running instance of YouTube Stats now. You can test this by opening `YOUR_IP_ADDRESS/ytstats/v1/` in your browser, and you should see some text indicating that you have reached the YouTube Stats REST API.

All you need to do now is to [get your YouTube API key](https://github.com/Travus/yt_stats/wiki#getting-a-youtube-api-key) and read up on what the different endpoints return in the [wiki](https://github.com/Travus/yt_stats/wiki) attacked to this repository.

## Contact
If you have any questions, needs or requests, feel free to contact me!  
I'm mostly active on Discord, but you can reach me on Twitter too.  
Alternatively, you can also create an issue on this repository and I will get back to you soon.

Discord: Travus#8888  
Twitter: [@RealTravus](https://twitter.com/RealTravus)
