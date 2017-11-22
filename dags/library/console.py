import sys


def write_tail(ctype):
    write_msg = "Finish Getting data[" + ctype + "]\r\n"
    sys.stdout.write(write_msg)
    sys.stdout.flush()


def write_head(ctype):
    write_msg = "Getting data[" + ctype + "]:\r\n"
    sys.stdout.write(write_msg)
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
