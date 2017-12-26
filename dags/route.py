from scripts import bash, test
import os
SCRIPT = os.environ.get("SCRIPT")


if SCRIPT == "test":
    test.start()
elif SCRIPT == "bash":
    bash.start()
