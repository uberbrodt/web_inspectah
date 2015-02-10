# web_inspectah
Http Fraud Detection for Virool challenge

This is the project that won the  DevWeek 2015 Hackathon challenge by Virool. The challenge was to create a server that processed web
requests and detected whether or not they were fraudulent. There were also performance requirements:
* All requests should respond within 100ms
* Server should be able to respond to 1000 requests a second

Full rules are here for now: http://devweek.virool.com/rules

I decided to use Go, which I've not used before but had been suggested to me for highly concurrent, fast
server applications. Virool provided training data with both good and bad sample requests, so there was a lot of opportunity for machine learning or other automated fraud detection techniques. My first approach however, was to isolate the very simple cases of fraud manually.

First, I collected all the 'good' user agents and used that as a whitelist for all requests user agents. Then I
disqualified any requests with localhost or private subnets in the ip address and referer fields (the idea is that
we're doing this for ad based services, so this makes sense). I then wrote a test runner (test_runner_virool.rb) to
see how my service performed with test data. I was only marking a little over a 1000 of the requests that were false as false so that was when I introduced the naive Bayesian classifier on the referrer email. At that point I figured I had about 85% accuracy on the training set (and was also out of time to complete it!). I would have liked to build out the bayesian part and do more work on finding a pattern in the IP addresses (perhaps doing reverse dns lookups?), but I'm proud of what I got done (since I started the challenge late on Saturday as I was working the Traitify booth most of the day).

In hindsight my only regret was not naming the project "Inspectah Web", since that's more like Inspectah Deck.
