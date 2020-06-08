console.log("test")
function validateForm() {
    let fonts = []
    let i
    for (i = 0; i < 6; i++) {
        let string = "font-" + i.toString()
        fonts.push(document.getElementById(string).value)
    }
    fonts.sort()
    for (i = 0; i < fonts.length - 1; i++) {
        if (fonts[i] == fonts[i+1]) {
            alert("Please make sure all entries are different")
            return false
        }
    }
    return true
}