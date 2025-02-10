return (hwInterfaces) => {
    console.log('hello world')
    hwInterfaces.display.clear()
    hwInterfaces.display.drawString('Hello :>', 0, 0)
    hwInterfaces.display.flip()
}