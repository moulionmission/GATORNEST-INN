import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const Register = ({ onToggle }) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    
    const [success, setSuccess] = useState('');
    const navigate = useNavigate(); // Initialize the navigate function
  
    const handleSubmit = async (e) => {
      e.preventDefault();
      setError('');
      setSuccess('');
  
      if (password !== confirmPassword) {
        setError('Passwords do not match');
        return;
      }
  
      try {
        const response = await axios.post('http://localhost:3000/register', {
          email,
          password
        });
  
        console.log('Registration successful:', response.data);
        setSuccess('Registration successful! You can now log in.');
      } catch (error) {
        console.error('Registration failed:', error.response?.data || error.message);
        setError(error.response?.data?.message || 'Registration failed. Please try again.');
      }
    };
  
  
    return (
      <div className="auth-container">
        <h2>Register</h2>
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
          <input 
            type="password" 
            placeholder="Confirm Password" 
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)} 
            required 
          />
          <button type="submit">Register</button>
        </form>
        <p>
          Already have an account?{' '}
          <span onClick={onToggle}>Login here</span>
        </p>
      </div>
    );
  };
  
  export default Register;