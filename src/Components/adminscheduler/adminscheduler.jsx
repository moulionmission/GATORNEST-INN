import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './AdminScheduler.css';

const roomTypes = [
  "Deluxe Suite", "Executive Room", "Standard Room", "Family Room",
  "Single Room", "King Suite", "Junior Suite", "VIP Room", "Budget Room"
];

const daysOfWeek = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
//const staffOptions = ['Alice', 'Bob', 'Charlie', 'Diana', 'Ethan', 'Fiona', 'George', 'Hannah'];
const timeSlots = ['Morning', 'Afternoon', 'Night'];

const StaffScheduler = () => {
  const [selectedRoomType, setSelectedRoomType] = useState(roomTypes[0]);
  const [roomCount, setRoomCount] = useState('');
  const [roomAvailability, setRoomAvailability] = useState({});
  const [selectedDay, setSelectedDay] = useState('Monday');
  const [selectedStaff, setSelectedStaff] = useState([]);
  const [selectedSlot, setSelectedSlot] = useState('Morning');
  const [schedule, setSchedule] = useState({});
  const [editing, setEditing] = useState({});
  const [staffOptions, setStaffOptions] = useState([]);
  const [reverseStaffMap, setReverseStaffMap] = useState({});

  useEffect(() => {
    const fetchSchedulesAndStaff = async () => {
      try {
        const [staffRes, scheduleRes] = await Promise.all([
          axios.get("http://localhost:3000/staffs", {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`, // If using JWT token for auth
                
              'Content-Type': 'application/json'
          }
          }),
          axios.get("http://localhost:3000/schedules", {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`, // If using JWT token for auth
                
              'Content-Type': 'application/json'
            }
          })
        ]);

        const staffList = staffRes.data;
        const scheduleList = scheduleRes.data;
        setStaffOptions(staffList.map(staff => staff.first_name));

        // Create map: staff_id -> first_name
        const staffMap = {};
        staffList.forEach(staff => {
          staffMap[staff.staff_id] = staff.first_name;
        });

        const reverseStaffMap = {};
        staffList.forEach(staff => {
          reverseStaffMap[staff.first_name] = staff.staff_id;
        });
        setReverseStaffMap(reverseStaffMap);


        console.log(staffMap);
        console.log(scheduleList);

        // Organize schedule by weekday
        const parsedSchedule = {};

        (scheduleList || []).forEach(entry => {
          const { staff_id, shift_date, shift_time } = entry;

          const dayOfWeek = shift_date;


          const name = staffMap[staff_id] || `Staff ${staff_id}`;
          if (!parsedSchedule[dayOfWeek]) parsedSchedule[dayOfWeek] = [];
          parsedSchedule[dayOfWeek].push({ name, slot: shift_time });

          console.log('Parsed:', { shift_date, dayOfWeek, name, slot: shift_time });
          console.log('Staff Name:', staffMap[staff_id]);

        });

        setSchedule(parsedSchedule);
        console.log("Final schedule state:", schedule);

      } catch (err) {
        console.error("Error fetching schedules or staff:", err);
      }
    };

    fetchSchedulesAndStaff();
  }, []);


  const handleRoomSave = (e) => {
    e.preventDefault();
    setRoomAvailability(prev => ({
      ...prev,
      [selectedRoomType]: parseInt(roomCount) || 0
    }));
    setRoomCount('');
  };

  const handleAssign = (e) => {
    e.preventDefault();
    if (!selectedStaff.length) return;

    setSchedule(prev => ({
      ...prev,
      [selectedDay]: [
        ...(prev[selectedDay] || []),
        ...selectedStaff.map(name => ({ name, slot: selectedSlot }))
      ]
    }));

    setSelectedStaff([]);
  };

  const startEditing = (day, index, entry) => {
    setEditing({
      day,
      index,
      name: entry.name,
      slot: entry.slot
    });
  };

  const saveEditedEntry = async () => {
    const { day, index, name, slot } = editing;
    const staff_id = reverseStaffMap[name];

    try {
      await axios.put(`http://localhost:3000/schedule/${staff_id}`, null, {
        params: {
          shift_date: day,
          new_shift_time: slot
        },
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        }
      });

      setSchedule(prev => {
        const updated = [...prev[day]];
        updated[index] = { name, slot };
        return { ...prev, [day]: updated };
    });

    setEditing({});
    } catch (err) {
      console.error("Failed to update schedule:", err);
      alert("Failed to update schedule on server.");
    }
  };

  const deleteEntry = async (day, index) => {
    const entry = schedule[day][index];
    const staff_id = reverseStaffMap[entry.name];
  
    try {
      await axios.delete(`http://localhost:3000/schedule/${staff_id}`, {
        params: {
          shift_date: day,
          shift_time: entry.slot,
        },
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
  
      // Only update frontend state after successful deletion
      setSchedule(prev => {
        const updated = [...prev[day]];
        updated.splice(index, 1);
        return { ...prev, [day]: updated };
      });
  
      if (editing.day === day && editing.index === index) {
        setEditing({});
      }
    } catch (err) {
      console.error("Failed to delete schedule:", err);
      alert("Failed to delete schedule from the server.");
    }
  };

  const assignStaff = async (day, slot, staffName) => {
    const staff_id = reverseStaffMap[staffName];
  
    try {
      const res = await axios.post("http://localhost:3000/schedule", {
        staff_id,
        shift_date: day,
        shift_time: slot,
      }, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
  
      // Update local state
      setSchedule(prev => ({
        ...prev,
        [day]: [...(prev[day] || []), { name: staffName, slot }]
      }));
  
      //alert("Staff assigned successfully!");
    } catch (err) {
      console.error("Failed to assign staff:", err);
      alert("Failed to assign staff to schedule.");
    }
  };
  

  const today = new Date().toLocaleString('en-US', { weekday: 'long' });
  console.log('day today:', today);
  const todayStaff = schedule[today] || [];

  return (
    <div className="staffScheduler">
      <h2>Hotel Admin Scheduler</h2>

      {/* Room Availability */}
      <div className="roomSection">
        <h3>Room Availability</h3>
        <form className="roomForm" onSubmit={handleRoomSave}>
          <div className="formGroup">
            <label>Select Room Type</label>
            <select value={selectedRoomType} onChange={e => setSelectedRoomType(e.target.value)}>
              {roomTypes.map(room => (
                <option key={room} value={room}>{room}</option>
              ))}
            </select>
          </div>

          <div className="formGroup">
            <label>Available Rooms</label>
            <input
              type="number"
              min="0"
              value={roomCount}
              onChange={e => setRoomCount(e.target.value)}
              placeholder="e.g. 5"
            />
          </div>

          <button type="submit" className="saveBtn">Save Availability</button>
        </form>

        <div className="availabilityList">
          <h4>Current Availability</h4>
          <ul>
            {Object.entries(roomAvailability).map(([room, count]) => (
              <li key={room}>{room}: {count} rooms</li>
            ))}
          </ul>
        </div>
      </div>

      {/* Staff Scheduling Form */}
      <form className="schedulerForm" onSubmit={handleAssign}>
        <h3>Assign Housekeeping Staff</h3>
        <div className="formGroup">
          <label>Select Day</label>
          <select value={selectedDay} onChange={e => setSelectedDay(e.target.value)}>
            {daysOfWeek.map(day => (
              <option key={day} value={day}>{day}</option>
            ))}
          </select>
        </div>

        <div className="formGroup">
          <label>Select Time Slot</label>
          <select value={selectedSlot} onChange={e => setSelectedSlot(e.target.value)}>
            {timeSlots.map(slot => (
              <option key={slot} value={slot}>{slot}</option>
            ))}
          </select>
        </div>

        <div className="formGroup">
          <label>Select Staff</label>
          <select
            multiple
            value={selectedStaff}
            onChange={(e) => setSelectedStaff(Array.from(e.target.selectedOptions, option => option.value))}
          >
            {staffOptions.map(staff => (
              <option key={staff} value={staff}>{staff}</option>
            ))}
          </select>
        </div>

        {/* <button type="submit" className="assignBtn">Assign Staff</button> */}
        <button 
          type="button"
          onClick={() => selectedStaff.forEach(name => assignStaff(selectedDay, selectedSlot, name))}
          className="assignBtn">Assign Staff
        </button>
      </form>

      {/* Weekly Schedule */}
      <div className="fullSchedule">
        <h3>Weekly Schedule</h3>
        {daysOfWeek.map(day => (
          <div key={day} className="dayBlock">
            <h4>{day}</h4>
            <ul>
              {(schedule[day] || []).length > 0 ? (
                schedule[day].map((entry, idx) => (
                  <li key={idx} className="editableEntry">
                    {editing.day === day && editing.index === idx ? (
                      <>
                        <select
                          value={editing.name}
                          onChange={e => setEditing(prev => ({ ...prev, name: e.target.value }))}
                        >
                          {staffOptions.map(name => (
                            <option key={name} value={name}>{name}</option>
                          ))}
                        </select>

                        <select
                          value={editing.slot}
                          onChange={e => setEditing(prev => ({ ...prev, slot: e.target.value }))}
                        >
                          {timeSlots.map(slot => (
                            <option key={slot} value={slot}>{slot}</option>
                          ))}
                        </select>

                        <button className="saveEdit" onClick={saveEditedEntry}>Save</button>
                      </>
                    ) : (
                      <>
                        {entry.name} – {entry.slot}
                        <div className="entryButtons">
                          <button className="editBtn" onClick={() => startEditing(day, idx, entry)}>Edit</button>
                          <button className="deleteBtn" onClick={() => deleteEntry(day, idx)}>Delete</button>
                        </div>
                      </>
                    )}
                  </li>
                ))
              ) : (
                <li className="no-staff">No staff assigned</li>
              )}
            </ul>
          </div>
        ))}
      </div>

      <div className="todaySection">
        <h3>Who's Working Today ({today})</h3>
        {todayStaff.length > 0 ? (
          <ul className="todayList">
            {todayStaff.map((entry, index) => (
              <li key={index}>{entry.name} – {entry.slot}</li>
            ))}
          </ul>
        ) : (
          <p>No staff scheduled for today.</p>
        )}
      </div>
    </div>
  );
};

export default StaffScheduler;
