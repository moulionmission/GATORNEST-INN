import { useState } from 'react';
import './Login.css';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const Login = ({ onToggle }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate(); // Initialize the navigate function

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(''); // Clear any previous error

    try {
      const response = await axios.post('http://localhost:3000/login', {
        email,
        password
      });

      console.log('Login successful:', response.data);
      localStorage.setItem('userEmail', email);
      localStorage.setItem('token', response.data.token); // Store JWT token
      navigate('/home'); // Redirect to homepage
    } catch (error) {
      console.error('Login failed:', error.response?.data || error.message);
      setError(error.response?.data?.message || 'Login failed. Please try again.');
    }
  };

  return (
    <div className="login-container">
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <input 
          type="email" 
          placeholder="Email" 
          value={email}
          onChange={(e) => setEmail(e.target.value)} 
          required 
        />
        <input 
          type="password" 
          placeholder="Password" 
          value={password}
          onChange={(e) => setPassword(e.target.value)} 
          required 
        />
        <button type="submit">Login</button>
      </form>
      <p>
        Don't have an account?{' '}
        <span onClick={onToggle} className="link">
          Register here
        </span>
      </p>
    </div>
  );
};

export default Login;
