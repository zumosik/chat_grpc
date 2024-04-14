import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

import {AuthServiceClient} from '@zumosik/grpc_chat_protos/js/auth/AuthServiceClientPb'
import { CreateUserRequest, CreateUserResponse } from '@zumosik/grpc_chat_protos/js/auth/auth_pb'

function App() {
  const [count, setCount] = useState(0)


  useEffect(() => { 
    (
      async () => {
        const client = new AuthServiceClient('http://localhost:7771')

        const req = new CreateUserRequest()
        req.setEmail("test_user_from_web_0@exmaple.org")
        req.setPassword("password")
        req.setUsername("test_user_from_web_0")
    
        const resp : CreateUserResponse = await client.createUser(req)
        console.log(resp.getSuccess())
        console.log(resp.getUser())
      }
    )()
  } , [])

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
