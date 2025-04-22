import React, { createContext, useState, useEffect } from 'react';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [isStaff, setIsStaff] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem('isStaff');
    setIsStaff(stored === 'true');
  }, []);

  return (
    <AuthContext.Provider value={{ isStaff, setIsStaff }}>
      {children}
    </AuthContext.Provider>
  );
};
