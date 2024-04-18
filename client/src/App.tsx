import './App.css'

import { createBrowserRouter, RouterProvider } from 'react-router-dom'

import { ThemeProvider } from '@/components/theme-provider'
import LoginPage from './pages/login'

import Home from './pages/home'

import { Toaster } from '@/components/ui/toaster'

function App() {

  const router = createBrowserRouter(
    [
      {
        path: '/',
        Component: () => {
          return (
            <Home/>
          )
        }
      }, {
        path: '/login',
        Component: () => {
          return <LoginPage isLogin={true} />
        }
      },
    ]
  )

  return (
    <>

      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <RouterProvider router={router} />
        <Toaster />
      </ThemeProvider>
    </>
  )
}

export default App
