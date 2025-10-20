import subprocess
import time
import os
import signal 
current_process = None

def playSound(filename):
    
    global current_process
    
    if current_process and current_process.poll() is None:
        try:
            current_process.terminate()
            try:
               current_process.wait(timeout=1)
            except subprocess.TimeoutExpired:
                print("proccess is being killed...")
                current_process.kill()
        except Exception as e:
            print(f"error: \n {e}")
    
    current_process = None
    
    path = "/home/marciano/Downloads"
    fullpath = os.path.join(path, filename)
    if not os.path.exists(fullpath):
        print("filepath does not exist")
        # need to download the file 
        # maybe default to a random downloaded file 
        return False

    try:
        process = subprocess.Popen(['cvlc', '--play-and-exit', fullpath],
                                    stdout=subprocess.PIPE,
                                    stderr=subprocess.PIPE)
        
        current_process = process
        return True
    except FileNotFoundError:
        print("ERROR: cvlc not found is VLC installed and in PATH")
    except Exception as e:
        print(f"ERROR: there is an error namely \n {e}")
        
        
def mmain():
    playSound("7NA.mp3")
    

#main()
