import subprocess
import time
import os

def playSound(filename):
    path = ""
    fullpath = os.path.join(path, filename)
    if not os.path.exists(fullpath):
        print("filepath does not exist")
        # need to download the file 
        # maybe default to a random downloaded file 
        return False

    try:
        proccess = subprocess.Popen(['cvlc', '--play-and-exit', fullpath],
                                    stdout=subprocess.PIPE,
                                    stderr=subprocess.PIPE)
        return True
    except FileNotFoundError:
        print("ERROR: cvlc not found is VLC installed and in PATH")
    except Exception as e:
        print(f"ERROR: there is an error namely \n {e}")