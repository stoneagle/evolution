from library import console

count_dict = dict()


def inc_by_index(index):
    if index in count_dict:
        count_dict[index] += 1
    else:
        count_dict[index] = 1
    return


def show_result():
    for index in count_dict:
        console.write_number(index, str(count_dict[index]))
        reset()
    return


def reset():
    global count_dict
    count_dict = dict()
    return
