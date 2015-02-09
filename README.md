# web_inspectah
Http Fraud Detection for Virool challenge

This is a project to solve the DevWeek challenge by Virool. The challenge was to create a server that processed web
requests and detected whether or not they were fraudulent. There were also performance requirements:
* All requests should respond within 100ms
* Server should be able to respond to 1000 requests a second

For this, I decided to use Go, which I've not used before but had been suggested to me for highly concurrent, fast
server applications. Virool provided training data with both good or bad data, so there was a lot of opportunity for
machine learning or other automated fraud detection techniques. My first approach however, was to isolate the very simple
cases of fraud manually.

First, I took all the 'false' user_agents and made a list to match against, failing others. Then I looked at the referrer
names and ip addresses. Anything localhost or on a private subnet was also marked as false. I then devloped a test
runner to see how my serveice was doing. I was only marking a little over a 1000 of the requests that were false as false
so that was when I introduced the naive Bayesian classifier on the referrer email. At that point I figured I had about
>85&% accuracy on the training set.
