// This is a generic message class that allows any other object as the payload.
// The payload will be marshalled to JSON
export class Message<T> {
    channel: string;
    reply_channel: string;
    payload: T;
}
