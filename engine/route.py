from scripts import bash, ashare, bitmex, part
import sys
import getopt


def main(argv):
    script = ""
    func = ""
    try:
        # 这里的 h 就表示该选项无参数，i:表示 i 选项后需要有参数
        opts, args = getopt.getopt(argv, "hs:f:", ["script=", "func="])
    except getopt.GetoptError:
        print('Error: route.py -s <script> -f <func>')
        sys.exit(2)

    for opt, arg in opts:
        if opt == "-h":
            print('route.py -s <script> -f <func>')
            sys.exit()
        elif opt in ("-s", "--script"):
            script = arg
        elif opt in ("-f", "--func"):
            func = arg
    if script == "bitmex":
        if func == "test":
            bitmex.test()
        elif func == "watch":
            bitmex.monitor()
    elif script == "ashare":
        if func == "test":
            ashare.test()
        elif func == "watch":
            ashare.monitor()
    elif script == "part":
        if func == "test":
            part.test()
    elif script == "bash":
        bash.start()


if __name__ == "__main__":
    main(sys.argv[1:])
