import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { NoteService, Note } from '../../services/note-service';

@Component({
  selector: 'app-note-dashboard',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './note-dashboard.html',
  styleUrl: './note-dashboard.scss',
})
export class NoteDashboard implements OnInit {
  notes: Note[] = [];
  // TODO: Replace with actual user ID from authentication
  userId: string = 'testUser1'; 

  constructor(private noteService: NoteService, private router: Router) {}

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

  viewNote(noteId: string): void {
    this.router.navigate(['/notes', noteId]);
  }
}
