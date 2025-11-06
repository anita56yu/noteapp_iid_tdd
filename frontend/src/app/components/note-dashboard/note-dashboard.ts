import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NoteService, Note } from '../../services/note-service';

@Component({
  selector: 'app-note-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './note-dashboard.html',
  styleUrl: './note-dashboard.scss',
})
export class NoteDashboard implements OnInit {
  notes: Note[] = [];
  // TODO: Replace with actual user ID from authentication
  userId: string = 'testUser1'; 

  constructor(private noteService: NoteService) {}

  ngOnInit(): void {
    this.noteService.getAccessibleNotes(this.userId).subscribe({
      next: (notes) => {
        this.notes = notes;
      },
      error: (err) => {
        console.error('Error fetching notes', err);
      },
    });
  }
}
