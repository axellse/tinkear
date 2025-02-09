/*
tinkearOS boot manager version beta-1.0.0
*/

const hardwareConfigurations = {
    "m1" : {
        "display" : {
            "power" : D14, //(if applicable)
            "scl" : D0,
            "sda": D4 
        }
    }
}
const hardware = hardwareConfigurations["$HWConfig_ReplacedByInstaller$"]

//It is not recommended to touch the code below unless you know what you're doing.

const fs = require("Storage")
const SSD1306 = require("SSD1306")
const startTime = new Date()

const resources = {
    'osNotFound' : 'gECBAAAAAAAAAAAAAAAAAAAAAAAkCAAA44QABAAAAAAAAAB4IAgAARREAAQAAAAAAAAAzHXJMd0UBxjOOx3HOZwAAMwlKkpRE4SlJCSlKUpQAAAMJSx6URBEpSQkpSlL0AAAGCUqQlEURKUkJKUpShAAADAlKTHQ44cYxCSdJzmQAAAAAAAAAAAAAAAAAAAIAAAAMAAAAAAAAAAAAAAAMAAAAAB////////////////////+AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACEAABADABAAAAAAAAAAHvSvAe98Al7wHvSAAAAgABCUqQEpEAdSkAKUgAAAcAAQlKkBKRACUpAelIAAAHAAHvevAS8QAlLwHpeAAAD4AAAAAAAAAAAAAAAAgAAB/AAAAEAoAgPvgCIAEAAAAdwAF5H3q94CKAAikXwAAAPeABScQKvSAi+AIp0QAAAH3wAUhEeqEgIggBSFEAAAB98AFLxHq94D76AIvRAAAA/fgAAAAAAAAAAAAAAAAAAf38AEQABAIAQBAAAIAAAAH//AH3vAV6gHvXogPvAAAD/f4ARLwFS8BL1KIAiQAAA//+AESgBUpAShSqAIkAAAAAAABEvAVKQHvXvgCPAAAA/HgAAAAAAAAAAAAAAAAAAf78ACAACAAQBEAAQCAAIAHO4AHvReve8B97wF77yPgBhsABKUUqQpAES8BCI84gAYb4ASlVKl6QBEoAXiICIAGGfAHvfSve8ARLwF4j3iABhgwAAAAAAAAAAAAAAAAAAc4cAQAQgAMAQAgAAPvgAAH+/AHpV4DyAfXqe96KAAAA/HgBKVSAlwBFL3hQi+AAAAAAASlUgJIARSlD0IggAAAAAAHvV4DyAEUpe9D74AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP//////////////+AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGwAAAgAZ+xw8AAAAAAAAABsAAAIAGGA2ZgAAAAAAczObC2OE07xjY3AAAAAAAJnm2w22xbbYY2M8AAAAAAB4x9sNt8m32GNjDgAAAAAA2eYbDbYJthhjY0YAAAAAANs222220bbYYzZmAAAAAABvM5tts5DzjGMcPAAAAAAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    'bootManagerErr' : 'gECBAAAAAAAAAAAAAAAAAAAAAAAkCAAA44QABAAAAAAAAAB4IAgAARREAAQAAAAAAAAAzHXJMd0UBxjOOx3HOZwAAMwlKkpRE4SlJCSlKUpQAAAMJSx6URBEpSQkpSlL0AAAGCUqQlEURKUkJKUpShAAADAlKTHQ44cYxCSdJzmQAAAAAAAAAAAAAAAAAAAIAAAAMAAAAAAAAAAAAAAAMAAAAAB////////////////////+AAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAcAA+gACAAEAAAAAAAAAAAHAACPeA973wH3ve97wAAAD4AAiXgJSkQBUKQpegAAAB/AAIlACUpEARel6UIAAAAdwACJeA97xAEXpe96AAAAPeAAAAAAAAAAAAABAAAAAH3wAAAEIAACAAAQCAAABQB98APeleA974C9feongPUA/fgCEpUgJSIApRAq54AVAf38AhKVICUiAKUR6oQA9QH//APe9eAl4gClEer3gPUD/f4AAAAAAAAAAAAAAAAAA//+AgAAgAAAAAAAABAAAAAAAAPe96L3vAe973vUgPeA/HgCQoSiFDwEoSl6FICUAf78Al6EqvQgBKEpQhSAlAHO4AJeh770PAeh73oXgPQBhsAAAAAAAAAEAAgAAIAAAYb4AYBAEBAAIAAAiAAPvgGGfAE9XvAXveveA+94CKCBhgwDhV6QFIUqUgCJeAi+Ac4cAT1QkBS9KlIAiUAIggH+/AE9XvAXvepeAIl4D76A/HgAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAPoABAAiAAPgAAQAAAAAAACD3vUA+94CD3r0AAAAAAAAgl6HgCJeA+9CFAAAAAAAAIJQhIAiUAAoQvQAAAAAAAD6XvSAIl4D70L0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAMAAAAAAAAMAAACAAA973wCe8B973gLye997694JSgQB0oAVSh4Cl0oVCIpSCUoEAJKAEUoQApJKEXiKUg96BACegBF6HgKSehF4i9IIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    'osCrashed' : 'gECBAAAAAAAAAAAAAAAAAAAAAAAkCAAA44QABAAAAAAAAAB4IAgAARREAAQAAAAAAAAAzHXJMd0UBxjOOx3HOZwAAMwlKkpRE4SlJCSlKUpQAAAMJSx6URBEpSQkpSlL0AAAGCUqQlEURKUkJKUpShAAADAlKTHQ44cYxCSdJzmQAAAAAAAAAAAAAAAAAAAIAAAAMAAAAAAAAAAAAAAAMAAAAAB////////////////////+AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAA4AZgADh4MAAAAAAMABmAAOAGYABszDAAAAAADAAZgAHwD3jgxuA8ccDrcc8ceYAD+AZtsMZ4NpthrpttttmAA7gGbfDGHDZ5wYx5zb7ZAAe8Bm2Axow22GGM2G2w2QAPvgZtsGzMNtthrNttttgAD74DbOA4eDZtwOxtzZx5gB+/AAAAAAAAAAAAAAAAAAA/v4AAAAAAAAAAAAAAAAAAP/+H////////////////gH+/wAAAAAAAAAAAAAAAAAB//8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB+PBECAAAAAAgAAAAgAAAA/34RevIQBL0j3gPe+AAAAOdwEUrznASlIh4CUiAAAADDYBFKgIQEpSIQAlIgAAAAw3wfSve8B73iHgJeIAAAAMM+AAAAAACAAAAAAAAAAADDBgAAAABACAAA++AAAgAA5w4SRe8B9ep73ooADyKAAP9+EnUpAEUveFCL4AE7wAB+PBIVKQBFKUPQiCAPCkAAAAAe9S8ARSl70PvoD3pAAAAAAAABAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAAEAQAAAAIgAAAAAAAAACe8B714D3gPveAAAAAAAAB0oAS9SAlIAiXgAAAAAAAAJKAEoUgJSAIlAAAAAAAAACegBL14D0gCJeAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAGIAAAIAAAAAAAAnvS+A94BK8B4D0vAAAAAAdKEqgJQA6vACAlKQAAAAACShKICUAEqAHgJSkAAAAAAnoeiA9ABK8B4D3vAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAe973vgAAAAAAAAAAAAAAAEPSlAgAAAAAAAAAAAAAAABCEpQIAAAAAAAAAAAAAAAAQ970CIAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA='
}

let serialInterface = {
    "installOs" : (name, code) => {
        try {
            code = atob(code)
            let osList = fs.read('os.json');
            if (osList) {
                osList = JSON.parse(osList)
            } else {
                osList = []
            }

            if (osList.length > 0) {
                console.log(JSON.stringify({
                    "status" : "error",
                    "message" : "Cant load multiple OS"
                }))
            }

            osList.push({
                "name" : name,
                "src" : code
            })
            fs.write('os.json', JSON.stringify(osList))
        } catch (e) {
            console.log(JSON.stringify({
                "status" : "error",
                "message" : "Something went wrong, either the os you tried to load or your os.json file may be corrupted.",
                "errDump" : e
            }))
        }
    },
    "clearOS" : () => {
        console.log('Alright, now clearing/uninstalling tinkearOS.')
        console.log('NOTE: User data will still be kept.')

        fs.erase('os.json')
    },
    "clearDevice" : () => {
        console.log('Completly clearing the device including all user data and all installed versions of tinkearOS')

        fs.eraseAll()
    }
} 

let hwInterfaces = {}
function showBootError(type, goal, fix, e) {
    if (!e) {
        e = 'N/A'
    }
    if (type == 'err') {
        console.log('The boot manager encountered a fatal error.')
        console.log('--------------------------------------------------')
        console.log('Encountered an error trying to ' + goal + ', try ' + fix + '.\nError dump:\n\n' + e)
        console.log('--------------------------------------------------')
        if (hwInterfaces.display) {
            hwInterfaces.display.drawImage(atob(resources.bootManagerErr))
            hwInterfaces.display.drawString('cant ' + goal, 1, 38)
            hwInterfaces.display.flip()
        }
    } else if (type == 'noOs') {
        console.log('No OS Loaded.')
        hwInterfaces.display.drawImage(atob(resources.osNotFound))
        hwInterfaces.display.flip()
    } else if (type == 'osCrash') {
        console.log('The operating system crashed.')
        console.log('--------------------------------------------------')
        console.log(e)
        console.log('--------------------------------------------------')
        hwInterfaces.display.drawImage(atob(resources.osCrashed))
        hwInterfaces.display.flip()
    }
}

function tryPoint(func, goal, fix) {
    try {
        return func()
    } catch (e) {
        showBootError('err', goal, fix, e)
        return false
    }
}

let _callOnce = false
function boot() {
    if (_callOnce) return
    _callOnce = true

    let osList = fs.read('os.json')
    if (!osList) {
        showBootError('noOs')
        return
    }
    if (!tryPoint(() => osList = JSON.parse(osList), 'parse os.json', 'fixing syntax errors in os.json')) return;

    if (!Array.isArray(osList)) {
        showBootError('err', 'parse os.json', 'making os.json an array')
    }

    if (osList.length > 1) {
        showBootError('err', "load multiple os", 'having one item in os.json')
    }

    //TODO: digitally sign official builds of tinkearOS with rsa and warn the user if they're about to load an inofficial build. (requires custom build of espruino)
    const launchOS = tryPoint(() => {return eval(osList[0].src)}, 'load the os', "asking for help on the forum or filing a bug report as long as you're using an offical build of tinkearOS.")
    try {
        console.log('Alright, attempting to launch "' + osList[0].name + '"')
        console.log('Took ' + (Date.now() - startTime) /1000 + 's for bootloader to initalize hardware and launch the OS.')
        launchOS(hwInterfaces)
    } catch (e) {
        showBootError('osCrash', undefined, undefined, e)
    }
}
tryPoint(() => {
    if (hardware.display.power) {
        digitalWrite(hardware.display.power, 1)
    }
    I2C1.setup({
        scl: hardware.display.scl,
        sda: hardware.display.sda
    })
    hwInterfaces.display = SSD1306.connect(I2C1, boot)
}, 'init the display', 'changing the hardware config')
