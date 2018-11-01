import matplotlib.pyplot as plt
import numpy as np
import json
from apiclient import APIClient

class PublicAPI(APIClient):
    BASE_URL = 'http://localhost:8001'

api = PublicAPI()

resp = api.call("/original")
print(resp)
respJson = json.loads(resp)

xs = np.array(respJson["xs"])
ys = np.array(respJson["ys"])

xi = np.arange(0, 9)
yi = np.sin(2 * np.pi * xi)

fig, ax = plt.subplots()
ax.plot(ys, xs, 'o', yi, xi)
ax.grid(True, linestyle='-')
ax.tick_params(labelcolor='b', labelsize='medium', width=1)
plt.show()