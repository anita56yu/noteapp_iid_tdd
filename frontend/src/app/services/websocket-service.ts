import { Injectable } from '@angular/core';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {
  private socket$: WebSocketSubject<any> | null = null;

  constructor() { }

  connect(noteId: string): Observable<any> {
    if (!this.socket$ || this.socket$.closed) {
      this.socket$ = webSocket(`ws://localhost:8080/notes/${noteId}/ws`);
    }
    return this.socket$.asObservable();
  }

  disconnect(): void {
    if (this.socket$) {
      this.socket$.complete();
      this.socket$ = null;
    }
  }

  sendMessage(message: any): void {
    if (this.socket$) {
      this.socket$.next(message);
    }
  }
}
