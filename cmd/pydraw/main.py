import matplotlib.pyplot as plt
import numpy as np
import json

# XStart, XEnd, YStart, YEnd

iter_max = 2
intensity = 0
format = "./cache/intensity-{}-{}.json"


fig = plt.figure()
axs = (fig.add_subplot(121, projection='3d'), fig.add_subplot(122, projection='3d'))
colors = ("blue", "green")

for iter in range(iter_max):
    with open(format.format(iter, intensity)) as json_file:
        dataDRT = json.load(json_file)

    data = {}

    for key in dataDRT:
        splited = key[1:-1].split(",")
        data[(float(splited[0]), float(splited[1]), float(splited[2]), float(splited[3]))] = dataDRT[key]

    xStarts = []
    yStarts = []
    zStarts = np.zeros(len(data))

    xEnds = []
    yEnds = []
    zEnds = []

    for reg in data:
        xStarts.append(reg[0])
        yStarts.append(reg[2])
        xEnds.append(reg[1] - reg[0])
        yEnds.append(reg[3] - reg[2])

        zEnds.append(data[reg])

    axs[iter].bar3d(xStarts, yStarts, zStarts, xEnds, yEnds, zEnds, color=colors[iter])

plt.show()