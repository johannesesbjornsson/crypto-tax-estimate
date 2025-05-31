import { useState } from 'react'
import './App.css'
import TransactionsPage from './components/TransactionsPage';

import Button from './components/Button';
import Tabs from './components/Tabs';
import UserProfile from './components/UserProfile';


function App() {
  const tabs = {
    Overview: <div className="">Overview content</div>,
        Transactions: <TransactionsPage />,
        Profile: (
      <UserProfile/>
    ),
  };

  return (
    <div className="">
      <Tabs tabs={tabs} />
    </div>
  );
}

export default App;
