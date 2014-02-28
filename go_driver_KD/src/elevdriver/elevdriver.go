package elevdriver
import "fmt"
import "time"

type Direction int

const (
    NONE Direction = iota
    UP
    DOWN
)
var prevDir Direction = NONE

type LightVal int
const (
    ON LightVal = iota
    OFF
)

type Button struct {
    Floor int
    Dir Direction
}


const MAX_SPEED = 4024
const MIN_SPEED = 2048

func Init(  buttonEventChan         chan Button,
            floorEventChan          chan int,
            stopButtonEventChan     chan bool,
            obstructionEventChan    chan bool){

    val := IoInit()
    if !val {
        fmt.Printf("Driver initiated\n")
    } else {
        fmt.Printf("Driver not initiated\n")
    }


    go poller(  buttonEventChan,
                floorEventChan,
                stopButtonEventChan,
                obstructionEventChan)
}

func poller(buttonEventChan         chan Button,
            floorEventChan          chan int,
            stopButtonEventChan     chan bool,
            obstructionEventChan    chan bool) {

    var floorMap = map[int] int {
        SENSOR1 : 0,
        SENSOR2 : 1,
        SENSOR3 : 2,
        SENSOR4 : 3,
    }

    var buttonMap = map[int] Button {
        FLOOR_COMMAND1 : { 1, NONE },
        FLOOR_COMMAND2 : { 2, NONE },
        FLOOR_COMMAND3 : { 3, NONE },
        FLOOR_COMMAND4 : { 4, NONE },
        FLOOR_UP1      : { 1,   UP },
        FLOOR_UP2      : { 2,   UP },
        FLOOR_UP3      : { 3,   UP },
        FLOOR_DOWN2    : { 2, DOWN },
        FLOOR_DOWN3    : { 3, DOWN },
        FLOOR_DOWN4    : { 4, DOWN },
    }

    buttonList := make(map[int]bool)
    for key, _ := range buttonMap {
        buttonList[key] = Read_bit(key)
    }

    floorList := make(map[int]bool)
    for key, _ := range floorMap {
        floorList[key] = Read_bit(key)
    }

    oldStop := false
    oldObs := false

    for {
        time.Sleep(50*time.Millisecond)
        
        for key, floor := range floorMap {
            newValue := Read_bit(key)
            if newValue != floorList[key] {
                newFloor := floor
                floorEventChan <- newFloor
            }
            floorList[key] = newValue
        }

        for key, btn := range buttonMap {
            newValue := Read_bit(key)
            if newValue && !buttonList[key] {
                newButton := btn
                buttonEventChan <- newButton
            }
            buttonList[key] = newValue
        }

        newStop := Read_bit(STOP)
        if newStop && !oldStop {
            stopButtonEventChan <- true
        }
        oldStop = newStop

        newObs := Read_bit(OBSTRUCTION)
        if newObs != oldObs {
            obstructionEventChan <- newObs
        }
        oldObs = newObs
    }

}


func SetMotorDir(newDir Direction) {
    if (newDir == NONE) && (prevDir == UP) {
        Set_bit(MOTORDIR)
        Write_analog(MOTOR, MIN_SPEED)
    } else if (newDir == NONE) && (prevDir == DOWN) {
        Clear_bit(MOTORDIR)
        Write_analog(MOTOR, MIN_SPEED)
    } else if (newDir == UP) {
        Clear_bit(MOTORDIR)
        Write_analog(MOTOR, MAX_SPEED)
    } else if (newDir == DOWN) {
        Set_bit(MOTORDIR)
        Write_analog(MOTOR, MAX_SPEED)
    } else {
        Write_analog(MOTOR, MIN_SPEED)
    }
    prevDir = newDir
}


// Fugly. Use arrays next time plz
// SetOrderButtonLight()? Stop is also a button...
func SetButtonLight(floor int, dir Direction, onoff LightVal) {
    switch onoff {
    case ON:
        switch {
        case    floor == 1 && dir == NONE:
                Set_bit(LIGHT_COMMAND1)
        case    floor == 2 && dir == NONE:
                Set_bit(LIGHT_COMMAND2)
        case    floor == 3 && dir == NONE:
                Set_bit(LIGHT_COMMAND3)
        case    floor == 4 && dir == NONE:
                Set_bit(LIGHT_COMMAND4)
        case    floor == 1 && dir == UP:
                Set_bit(LIGHT_UP1)
        case    floor == 2 && dir == UP:
                Set_bit(LIGHT_UP2)
        case    floor == 3 && dir == UP:
                Set_bit(LIGHT_UP3)
        case    floor == 2 && dir == DOWN:
                Set_bit(LIGHT_DOWN2)
        case    floor == 3 && dir == DOWN:
                Set_bit(LIGHT_DOWN3)
        case    floor == 4 && dir == DOWN:
                Set_bit(LIGHT_DOWN4)
        }
    case OFF:
        switch {
        case    floor == 1 && dir == NONE:
                Clear_bit(LIGHT_COMMAND1)
        case    floor == 2 && dir == NONE:
                Clear_bit(LIGHT_COMMAND2)
        case    floor == 3 && dir == NONE:
                Clear_bit(LIGHT_COMMAND3)
        case    floor == 4 && dir == NONE:
                Clear_bit(LIGHT_COMMAND4)
        case    floor == 1 && dir == UP:
                Clear_bit(LIGHT_UP1)
        case    floor == 2 && dir == UP:
                Clear_bit(LIGHT_UP2)
        case    floor == 3 && dir == UP:
                Clear_bit(LIGHT_UP3)
        case    floor == 2 && dir == DOWN:
                Clear_bit(LIGHT_DOWN2)
        case    floor == 3 && dir == DOWN:
                Clear_bit(LIGHT_DOWN3)
        case    floor == 4 && dir == DOWN:
                Clear_bit(LIGHT_DOWN4)
        }
    }
}


func SetFloorLight(floor int) {
    switch floor {
    case 1:
        Clear_bit (FLOOR_IND1)
        Clear_bit (FLOOR_IND2)
    case 2:
        Clear_bit (FLOOR_IND1)
        Set_bit   (FLOOR_IND2)
    case 3:
        Set_bit   (FLOOR_IND1)
        Clear_bit (FLOOR_IND2)
    case 4:
        Set_bit   (FLOOR_IND1)
        Set_bit   (FLOOR_IND2)
    }
}


func SetStopButtonLight(onoff LightVal) {
    switch {
    case onoff == ON:
        Set_bit(LIGHT_STOP)
    case onoff == OFF:
        Clear_bit(LIGHT_STOP)
    }
}


func SetDoorOpenLight(onoff LightVal) {
    switch {
    case onoff == ON:
        Set_bit(DOOR_OPEN)
    case onoff == OFF:
        Clear_bit(DOOR_OPEN)
    }
}