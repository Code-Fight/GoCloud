
//offiical : https://github.com/qiniu/qetag/


// ref: https://www.jianshu.com/p/3785fc314fc5

function getEtag(buffer, callback) {
    // sha1算法
    var shA1 =  sha1.digest;

    // 以4M为单位分割
    var blockSize = 4 * 1024 * 1024;
    var sha1String = [];
    var prefix = 0x16;
    var blockCount = 0;

    var bufferSize = buffer.size || buffer.length || buffer.byteLength;
    blockCount = Math.ceil(bufferSize / blockSize);

    for (var i = 0; i < blockCount; i++) {
        sha1String.push(shA1(buffer.slice(i * blockSize, (i + 1) * blockSize)));
    }
    function concatArr2Uint8(s) {//Array 2 Uint8Array
        var tmp = [];
        for (var i of s) tmp = tmp.concat(i);
        return new Uint8Array(tmp);
    }
    function Uint8ToBase64(u8Arr, urisafe) {//Uint8Array 2 Base64
        var CHUNK_SIZE = 0x8000; //arbitrary number
        var index = 0;
        var length = u8Arr.length;
        var result = '';
        var slice;
        while (index < length) {
            slice = u8Arr.subarray(index, Math.min(index + CHUNK_SIZE, length));
            result += String.fromCharCode.apply(null, slice);
            index += CHUNK_SIZE;
        }
        return urisafe ? btoa(result).replace(/\//g, '_').replace(/\+/g, '-') : btoa(result);
    }
    function calcEtag() {
        if (!sha1String.length) return 'Fto5o-5ea0sNMlW_75VgGJCv2AcJ';
        var sha1Buffer = concatArr2Uint8(sha1String);
        // 如果大于4M，则对各个块的sha1结果再次sha1
        if (blockCount > 1) {
            prefix = 0x96;
            sha1Buffer = shA1(sha1Buffer.buffer);
        } else {
            sha1Buffer = Array.apply([], sha1Buffer);
        }
        sha1Buffer = concatArr2Uint8([[prefix], sha1Buffer]);
        return Uint8ToBase64(sha1Buffer, true);
    }
    return (calcEtag());
}