import sys


head_action = ""
head_resource = ""
head_detail = ""


def write_tail():
    global head_action
    global head_resource
    global head_detail
    write_msg = "Finish " + head_action + " " + head_resource + "[" + head_detail + "]\r\n"
    sys.stdout.write(write_msg)
    sys.stdout.flush()


def write_head(action, resource, detail):
    global head_action
    global head_resource
    global head_detail
    head_action = action
    head_resource = resource
    head_detail = detail
    write_msg = "Execing " + head_action + " " + head_resource + "[" + head_detail + "]\r\n"
    sys.stdout.write(write_msg)
    sys.stdout.flush()


def write_msg(msg):
    sys.stdout.write(msg + "\r\n")
    sys.stdout.flush()


def write_create():
    sys.stdout.write('C')
    sys.stdout.flush()


def write_blank():
    sys.stdout.write("\r\n")
    sys.stdout.flush()


def write_number(index, num):
    write_msg = "[" + index + " : " + num + "]\r\n"
    sys.stdout.write(write_msg)
    sys.stdout.flush()


def write_append():
    sys.stdout.write('A')
    sys.stdout.flush()


def write_pass():
    sys.stdout.write('P')
    sys.stdout.flush()


def write_error():
    sys.stdout.write('E')
    sys.stdout.flush()


def write_exec():
    sys.stdout.write('#')
    sys.stdout.flush()
