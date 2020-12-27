/**
 * Date:2020/12/24
 * Desc:
 */

export function getOneTimeRange(time, options) {
    let date;

    // 都为空的时候
    if (!time && !options) {
        date = new Date();
    } else if (Object.prototype.toString.call(time) !== '[object Date]' && time !== null && typeof time === 'object') {
        // time 为 options 参数。
        options = time;
        date = new Date();
    } else if (Object.prototype.toString.call(time) === '[object Date]') {
        // time 是时间格式
        date = time;
    } else {
        // time 是 int 格式。
        if (('' + time).length === 10) time = parseInt(time) * 1000;
        time = +time; // 转成int 型
        date = new Date(time);
    }

    options = options || {};

    let result = {
        start: 0,
        end: 0
    };

    let _startTime = new Date(date).setHours(options.startHour || 0, options.startMin || 0, 0, 0);
    let _endTime = new Date(date).setHours(options.endHour || 23, options.endMin || 59, 59, 0);
    result.start = new Date(_startTime).getTime();
    result.end = new Date(_endTime).getTime();

    return result;
};


export function formatTimestamp(t) {
    var d = 0,
        h = 0,
        m = 0,
        s = 0;
    if (t > 0) {
        d = Math.floor(t / 1000 / 3600 / 24)
        h = Math.floor(t / 1000 / 60 / 60 % 24)
        m = Math.floor(t / 1000 / 60 % 60)
        s = Math.floor(t / 1000 % 60)
    }

    return `${d}天${h}时${m}分${s}秒`
}

// 单位秒
export function formatTimeTips(timestamp) {
    let result;

    //
    if (timestamp > -1) {
        let hour = Math.floor(timestamp / 3600);
        let min = Math.floor(timestamp / 60) % 60;
        let sec = timestamp % 60;

        sec = Math.round(sec);

        if (hour < 10) {
            result = '0' + hour + ":";
        } else {
            result = hour + ":";
        }

        if (min < 10) {
            result += "0";
        }
        result += min + ":";
        if (sec < 10) {
            result += "0";
        }
        result += sec.toFixed(0);
    }

    return result;
}


export function parseTime(time, cFormat) {
    if (arguments.length === 0) {
        return null
    }
    var format = cFormat || '{y}-{m}-{d} {h}:{i}:{s}'
    var date;
    if (typeof time === 'object') {
        date = time
    } else {
        if (('' + time).length === 10) time = parseInt(time) * 1000;
        time = +time; // 转成int 型
        date = new Date(time)
    }
    var formatObj = {
        y: date.getFullYear(),
        m: date.getMonth() + 1,
        d: date.getDate(),
        h: date.getHours(),
        i: date.getMinutes(),
        s: date.getSeconds(),
        a: date.getDay()
    };
    var time_str = format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
        var value = formatObj[key]
        if (key === 'a') return ['一', '二', '三', '四', '五', '六', '日'][value - 1]
        if (result.length > 0 && value < 10) {
            value = '0' + value
        }
        return value || 0
    });
    return time_str
}


export function isDef(v) {
    return v !== undefined && v !== null;
}
