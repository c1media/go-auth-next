"use client"

import { useState, useEffect, createContext, useContext } from 'react'
import React from 'react'
import type { User } from '@/auth'

export interface Session {
  user: User
  expires: string
}

export interface SessionContextValue {
  data: Session | null
  status: 'loading' | 'authenticated' | 'unauthenticated'
  update: () => Promise<void>
}

const SessionContext = createContext<SessionContextValue | undefined>(undefined)

export function SessionProvider({ children }: { children: React.ReactNode }) {
  const [session, setSession] = useState<Session | null>(null)
  const [status, setStatus] = useState<'loading' | 'authenticated' | 'unauthenticated'>('loading')

  const fetchSession = async () => {
    try {
      const response = await fetch('/api/auth/session')
      
      if (response.ok) {
        const data = await response.json()
        if (data.authenticated) {
          setSession({
            user: data.user,
            expires: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
          })
          setStatus('authenticated')
        } else {
          setSession(null)
          setStatus('unauthenticated')
        }
      } else {
        setSession(null)
        setStatus('unauthenticated')
      }
    } catch (error) {
      console.error('Session fetch error:', error)
      setSession(null)
      setStatus('unauthenticated')
    }
  }

  useEffect(() => {
    fetchSession()
  }, [])

  const contextValue: SessionContextValue = {
    data: session,
    status,
    update: fetchSession
  }

  return React.createElement(
    SessionContext.Provider,
    { value: contextValue },
    children
  )
}

export function useSession() {
  const context = useContext(SessionContext)
  if (context === undefined) {
    throw new Error('useSession must be used within a SessionProvider')
  }
  return context
}

// Utility function to get CSRF token from cookie
export function getCSRFToken(): string | null {
  if (typeof document === 'undefined') return null
  
  const cookie = document.cookie
    .split('; ')
    .find(row => row.startsWith('csrf='))
  
  return cookie ? cookie.split('=')[1] : null
}