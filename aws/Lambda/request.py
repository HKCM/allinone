import json
import urllib3


def lambda_handler(event, context):
    print(event)
    http = urllib3.PoolManager()
    url = event["url"]
    # "API4http":"http://xxxxxx/restapi/v1.0/account/~"
    # "API4https":"https://xxxxxx/restapi/v1.0/account/~"
    # "Media4http":"http://xxxxxx/restapi/v1.0/account/~/profile-image"
    # "Media4https":"https://xxxxxx/restapi/v1.0/account/~/profile-image"
    response = http.request('GET',url,headers = {'Content-Type': 'application/json'})
    print("Visit URL: " + url)
    print("Response Code: " + str(response.status))
    return response.status
    #response.data.decode()

# test event
# {"url": "http://www.baidu.com/restapi/v1.0/account/~"}