import subprocess
import time
import matplotlib.pyplot as plt

def runCommand(command):
    start_time = time.time()
    subprocess.run(command, check=True)
    end_time = time.time()
    return end_time - start_time

def getIdealSpeedUp(fileSize, threadNums=[2,4,6,8,12], address="/home/tianyushi/project-3-tianyushi/proj3/editor"):
    idealSpeedUp = {}
    sequentialTime = runCommand(["go", "run", address, fileSize])
    sequentialPartTime = runCommand(["go", "run", address, fileSize, "pipelines","-1"])
    parallelFraction = (sequentialTime-sequentialPartTime)/sequentialTime   

    for num in threadNums: 
        idealSpeedUp[num] = round(1 / ((1 - parallelFraction) + (parallelFraction / num)), 2)

    return idealSpeedUp

def getRealSpeedUp(fileSize, Mode, threadNums=[2,4,6,8,12], address="/home/tianyushi/project-3-tianyushi/proj3/editor"):
    realSpeedUp = {}
    sequentialTime = runCommand(["go", "run", address, fileSize])

    for num in threadNums: 
        parallelTime = runCommand(["go", "run", address, fileSize, Mode, str(num)])
        realSpeedUp[num] = round(sequentialTime / parallelTime, 2)

    return realSpeedUp

def drawGraph(data_dict, titlePrefix):
    plt.figure(figsize=(14, 7))

    for key, value in data_dict.items():
        thread_numbers = list(value.keys())
        speedups = list(value.values())
        plt.plot(thread_numbers, speedups, marker='o', label=key)

    plt.title(f"{titlePrefix} Speed Up Comparison")
    plt.xlabel("Number of Threads")
    plt.ylabel("Speed Up")
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(f"{titlePrefix}_Speed_Up_Comparison.png")
    plt.show()
    plt.close()

def performTest():
    sizes = ["small", "mixture", "big"]
    modes = ["pipelines", "workStealing"]
    threadNums = [2, 4, 6, 8, 12]

    ideal_data = {}
    real_data = {mode: {} for mode in modes}

    for size in sizes:
        ideal_data[f"Ideal SpeedUp ({size})"] = getIdealSpeedUp(size, threadNums)

        for mode in modes:
            real_data[mode][f"Real SpeedUp ({size}, {mode})"] = getRealSpeedUp(size, mode, threadNums)

    drawGraph(ideal_data, "Ideal")
    for mode, modeData in real_data.items():
        drawGraph(modeData, f"Real ({mode})")

def main(): 
    performTest()

if __name__ == "__main__":
    main()
