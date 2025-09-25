const express = require('express');
const cors = require('cors');
const app = express();
const PORT = 8080;

// Middleware to parse JSON and enable CORS
app.use(cors());
app.use(express.json());

// In-memory "database" for events
let events = [];
let nextId = 1;

// Simple Server-Sent Events implementation
let clients = [];
function sendNotification(data) {
    clients.forEach(client => client.write(`data: ${JSON.stringify(data)}\n\n`));
}

// --- Endpoints ---

// Create Event
app.post('/events', (req, res) => {
    const { title, description, date } = req.body;
    if (!title || !description || !date) {
        return res.status(400).send('Title, description and Date are required.');
    }
    console.log('Creating event:', title, description, date);
    const newEvent = { id: nextId++, title, description, date };
    events.push(newEvent);

    // Notify all clients about the new event
    sendNotification({
        type: 'new_event',
        ts: new Date().toISOString(),
        event: newEvent
    });

    res.status(200).send(newEvent);
});

// Get all Events
app.get('/events', (req, res) => {
    res.status(200).json(events);
});

// Update Event
app.patch('/events/:id', (req, res) => {
    const id = parseInt(req.params.id, 10);
    const { title, description } = req.body;

    const eventToUpdate = events.find(event => event.id === id);

    if (!eventToUpdate) {
        return res.status(404).send('Event not found.');
    }

    if (title) {
        eventToUpdate.title = title;
    }
    if (description) {
        eventToUpdate.description = description;
    }

    // Notify clients about update (optional)
    sendNotification({
        type: 'updated_event',
        ts: new Date().toISOString(),
        event: eventToUpdate
    });

    res.status(200).json(eventToUpdate);
});

// Delete Event
app.delete('/events/:id', (req, res) => {
    const id = parseInt(req.params.id, 10);
    const initialLength = events.length;

    // Filter out the event to be deleted
    events = events.filter(event => event.id !== id);

    if (events.length === initialLength) {
        return res.status(404).send('Event not found.');
    }

    // Notify clients about deletion
    sendNotification({
        type: 'deleted_event',
        ts: new Date().toISOString(),
        event: { id }
    });

    res.status(200).send('Event deleted successfully.');
});

// Start the server
app.listen(PORT, () => {
    console.log(`Server listening on port ${PORT}`);
});
