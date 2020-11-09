/**
 * Date:2020/11/9
 * Desc:
 */
/**
 * Date:2020/11/2
 * Desc: ptz cmd 封装
 * cmd[0] //首字节以05H开头
 * cmd[1] //组合码，高4位为版本信息v1.0,版本信息0H，低四位为校验码
 *        //  校验码 = (cmd[0]的高4位+cmd[0]的低4位+cmd[1]的高4位)%16
 * cmd[2] //地址的低8位？？？什么地址，地址范围000h ~ FFFh(0~4095),其中000h为广播地址
 * cmd[3] //指令码
 * cmd[4] //数据1,水平控制速度、聚焦速度
 * cmd[5] //数据2，垂直控制速度、光圈速度
 * cmd[6] // 高4位为数据3=变倍控制速度，低4位为地址高4位
 */

const PTZ_TYPE = {
    stop: 'stop',
    right: 'right',
    left: 'left',
    up: 'up',
    down: 'down',
    leftUp: 'leftUp',
    leftDown: 'leftDown',
    rightUp: 'rightUp',
    rightDown: 'rightDown',
    zoomFar: 'zoomFar',
    zoomNear: 'zoomNear',
    apertureFar: 'apertureFar',
    apertureNear: 'apertureNear',
    focusFar: 'focusFar',
    focusNear: 'focusNear',
    setPos: 'setPos',
    calPos: 'calPos',
    delPos: 'delPos'
};

const PTZ_CMD_TYPE = {
    stop: 0x00,

    right: 0x01,
    left: 0x02,
    up: 0x08,
    down: 0x04,

    leftUp: 0x0A,
    leftDown: 0x06,
    rightUp: 0x09,
    rightDown: 0x05,

    zoomFar: 0x10, // 镜头 放大
    zoomNear: 0x20,  // 镜头 缩小

    apertureFar: 0x48, // 光圈 缩小
    apertureNear: 0x44, // 光圈 放大

    focusFar: 0x42, // 聚焦 近
    focusNear: 0x41, // 聚焦 远

    setPos: 0x81,
    calPos: 0x82,
    delPos: 0x83
};

const SPEED_ARRAY = [0x19, 0x32, 0x4b, 0x64, 0x7d, 0x96, 0xAF, 0xC8, 0xE1, 0xFA];
const POSITION_ARRAY = [0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10];

// 获取 direction 方向型
/**
 *
 * @param options
 *        type:
 *        speed:default 5
 *        index:
 * @returns {string}
 */
function getPTZCmd(options) {
    const {type, speed, index} = options;
    const ptzSpeed = getPTZSpeed(speed);
    let indexValue3, indexValue4, indexValue5, indexValue6;
    indexValue3 = PTZ_CMD_TYPE[type];
    switch (type) {
        case PTZ_TYPE.up:
        case PTZ_TYPE.down:
        case PTZ_TYPE.apertureFar:
        case PTZ_TYPE.apertureNear:
            indexValue5 = ptzSpeed;
            break;
        case PTZ_TYPE.right:
        case PTZ_TYPE.left:
        case PTZ_TYPE.focusFar:
        case PTZ_TYPE.focusNear:
            indexValue4 = ptzSpeed;
            break;
        case PTZ_TYPE.leftUp:
        case PTZ_TYPE.leftDown:
        case PTZ_TYPE.rightUp:
        case PTZ_TYPE.rightDown:
            indexValue4 = ptzSpeed;
            indexValue5 = ptzSpeed;
            break;
        case PTZ_TYPE.zoomFar:
        case PTZ_TYPE.zoomNear:
            indexValue6 = 0x10;
            break;
        case PTZ_TYPE.calPos:
        case PTZ_TYPE.delPos:
        case PTZ_TYPE.setPos:
            indexValue5 = getPTZPositionIndex(index);
            break;
        default:
            break;
    }
    return ptzCmdToString(indexValue3, indexValue4, indexValue5, indexValue6);
}

function getPTZSpeed(speed) {
    speed = speed || 5;
    const speedIndex = speed - 1;
    const ptzSpeed = SPEED_ARRAY[speedIndex] || SPEED_ARRAY[4];
    return ptzSpeed;
}

function getPTZPositionIndex(index) {
    return POSITION_ARRAY[index - 1];
}

function ptzCmdToString(indexValue3, indexValue4, indexValue5, indexValue6) {
    //
    let cmd = Buffer.alloc(8);
    // 首字节以05H开头
    cmd[0] = 0xA5;
    // 组合码，高4位为版本信息v1.0,版本信息0H，低四位为校验码
    cmd[1] = 0x0F;
    // 校验码 = (cmd[0]的高4位+cmd[0]的低4位+cmd[1]的高4位)%16
    cmd[2] = 0x01;
    //
    if (indexValue3) {
        cmd[3] = indexValue3;
    }
    if (indexValue4) {
        cmd[4] = indexValue4;
    }
    if (indexValue5) {
        cmd[5] = indexValue5;
    }
    if (indexValue6) {
        cmd[6] = indexValue6;
    }

    cmd[7] = (cmd[0] + cmd[1] + cmd[2] + cmd[3] + cmd[4] + cmd[5] + cmd[6]) % 256;

    return bytes2HexString(cmd);
}

function bytes2HexString(byte) {
    let hexs = "";
    for (let i = 0; i < byte.length; i++) {
        let hex = (byte[i]).toString(16);
        if (hex.length === 1) {
            hex = '0' + hex;
        }
        hexs += hex.toUpperCase();
    }
    return hexs;
}

export {
    getPTZCmd,
    PTZ_TYPE
}
