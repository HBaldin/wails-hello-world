import logo from './assets/images/logo-universal.png'
import {GreetingPage} from './features/greeting'
import {UpdatePanel} from './features/update'
import './App.css'

function App() {
    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo" />
            <GreetingPage />
            <UpdatePanel />
        </div>
    )
}

export default App
