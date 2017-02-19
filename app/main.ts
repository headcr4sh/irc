import { app, BrowserWindow } from "electron";

let mainWindow: Electron.BrowserWindow | null;

// Quit when all windows are closed.
app.on("window-all-closed", () => {
    if (process.platform != "darwin")
        app.quit();
});

// This method will be called when Electron has done everything
// initialization and ready for creating browser windows.
app.on("ready", () => {
    // Create the browser window.
    mainWindow = new BrowserWindow({
        width: 800, height: 600
    });

    // and load the index.html of the app.
    mainWindow.loadURL(`file://${__dirname}/www/index.html`);

    // Emitted when the window is closed.
    mainWindow.on("closed", () => {
        // Dereference the window object, usually you would store windows
        // in an array if your app supports multi windows, this is the time
        // when you should delete the corresponding element.
        mainWindow = null;
    });
});
