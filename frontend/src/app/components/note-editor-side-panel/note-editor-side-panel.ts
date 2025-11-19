import { Component, Input, Output, EventEmitter, OnChanges, SimpleChanges, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NoteService, Note, Content } from '../../services/note-service';
import { WebSocketService } from '../../services/websocket-service';
import { merge, Subscription } from 'rxjs';

@Component({
  selector: 'app-note-editor-side-panel',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './note-editor-side-panel.html',
  styleUrl: './note-editor-side-panel.scss',
})
export class NoteEditorSidePanelComponent implements OnChanges, OnDestroy {
  @Input() noteId: string | null = null;
  @Output() closePanel = new EventEmitter<void>();

  note: Note | null = null;
  private wsSubscription: Subscription | null = null;

  constructor(private noteService: NoteService, private webSocketService: WebSocketService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['noteId'] && this.noteId) {
      this.noteService.getNoteById(this.noteId).subscribe({
        next: (note) => {
          this.note = note;
          console.log('Loaded note', note);

          if (!note.contents || note.contents.length === 0) {
            const newContent: Content = {
              id: '', // ID will be set by the backend
              noteId: note.id,
              data: '',
              type: 'text',
              version: 0,
              position: 0,
            };
            this.noteService.addContent(note.id, newContent, note.version).subscribe({
              next: (addedContent) => {
                if (this.note) {
                  const createdContent: Content = { ...newContent, id: addedContent.id, version: 0 };
                  this.note.contents = [createdContent];
                  this.note.version++; // Increment note version as content was added
                }
              },
              error: (err) => {
                console.error('Error adding initial content to new note', err);
              },
            });
          }

          this.webSocketService.disconnect();
          this.wsSubscription = this.webSocketService.connect(note.id).subscribe({
            next: (message) => this.handleWebSocketMessage(message),
            error: (err) => console.error('WebSocket error:', err),
            complete: () => console.log('WebSocket connection closed'),
          });
        },
        error: (err) => {
          console.error('Error fetching note', err);
          this.note = null;
        },
      });
    }
  }

  ngOnDestroy(): void {
    this.webSocketService.disconnect();
    if (this.wsSubscription) {
      this.wsSubscription.unsubscribe();
    }
  }

  private handleWebSocketMessage(message: any): void {
    console.log('Received WebSocket message:', message);
    if (!this.note || message.note_id !== this.note.id) {
      return; // Ignore messages for other notes or if note is not loaded
    }

    // Update note version if the incoming message has a newer version
    if (message.type !== 'update_content' && message.note_version <= this.note.version) {
      return;
    }

    switch (message.type) {
      case 'add_content':
        const newContent: Content = {
          id: message.content_id,
          noteId: message.note_id,
          data: message.data,
          type: message.content_type,
          version: message.content_version,
          position: message.index,
        };
        this.note.contents.splice(message.index, 0, newContent);
        this.note.version = message.note_version;
        break;
      case 'update_content':
        const contentToUpdate = this.note.contents.find(c => c.id === message.content_id);
        var paragraph = document.querySelector(`[data-content-id="${message.content_id}"]`) as HTMLElement;
        console.log(contentToUpdate);
        if (contentToUpdate && message.content_version > contentToUpdate.version && message.data !== contentToUpdate.data) {
          console.log('original Content:', contentToUpdate);
          contentToUpdate.data = message.data;
          paragraph.textContent = message.data;
          console.log('updated Content:', contentToUpdate);
          contentToUpdate.version = message.content_version;
        }
        break;
      case 'delete_content':
        const contentIndexToDelete = this.note.contents.findIndex(c => c.id === message.content_id);
        if (contentIndexToDelete !== -1) {
          this.note.contents.splice(contentIndexToDelete, 1);
        }
        this.note.version = message.note_version;
        break;
      case 'update_note':
        console.log('Note title updating via WebSocket:', this.note.title, '->', message.data);

        if (message.data && message.data !== this.note.title) {
          this.note.title = message.data;
        }
        this.note.version = message.note_version;
        break;
      // case 'delete_note':
      //   this.onClose();
      //   break;
      default:
        console.warn('Unknown WebSocket message type:', message.type);
    }
  }

  onClose(): void {
    this.closePanel.emit();
  }

  onTitleBlur(event: FocusEvent): void {
    const newTitle = (event.target as HTMLElement).textContent?.trim() || '';
    if (this.note && newTitle !== this.note.title) {
      this.noteService.updateNote(this.note.id, newTitle, this.note.version).subscribe({
        next: () => {
          console.log('Note title updated successfully');
        },
        error: (err) => {
          console.error('Error updating note title', err);
          // Revert the UI change on error
          if (event.target) {
            (event.target as HTMLElement).textContent = this.note?.title ? this.note.title : '';
          }
        },
      });
    }
  }

  onTitleKeydownEnter(event: Event): void {
    const keyboardEvent = event as KeyboardEvent;
    keyboardEvent.preventDefault(); // Prevent new line
    (keyboardEvent.target as HTMLElement).blur(); // Trigger blur to save changes
  }

  onContentBlur(event: FocusEvent): void {
    console.log('onContentBlur called');
    setTimeout(() => {
      const paragraph = event.target as HTMLElement;
      const contentId = paragraph.getAttribute('data-content-id');
      if (contentId) {
        this.updateContent(contentId, paragraph.textContent);
      }
    }, 0);
  }

  onContentKeydown(event: KeyboardEvent): void {
    const paragraph = event.target as HTMLElement;
    const contentId = paragraph.getAttribute('data-content-id');

    if (event.key === 'ArrowUp' || event.key === 'ArrowDown') {
      event.preventDefault();
      if (!contentId || !this.note) return;

      const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
      if (currentIndex === -1) return;

      let targetIndex = -1;
      if (event.key === 'ArrowUp' && currentIndex > 0) {
        targetIndex = currentIndex - 1;
      } else if (event.key === 'ArrowDown' && currentIndex < this.note.contents.length - 1) {
        targetIndex = currentIndex + 1;
      }

      if (targetIndex !== -1) {
        const targetContentId = this.note.contents[targetIndex].id;
        const targetElement = document.querySelector(`[data-content-id="${targetContentId}"]`) as HTMLElement;
        if (targetElement) {
          targetElement.focus();
        }
      }
      return;
    }

    if (event.key === 'Enter') {
      event.preventDefault();
      if (!contentId || !this.note) return;

      const selection = window.getSelection();
      if (!selection || !selection.rangeCount) return;

      const range = selection.getRangeAt(0);
      const { startOffset } = range;

      const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
      if (currentIndex === -1) return;

      const originalText = paragraph.textContent || '';
      const textBeforeCursor = originalText.substring(0, startOffset);
      const textAfterCursor = originalText.substring(startOffset);

      // Update current content in the DOM immediately for responsiveness
      paragraph.textContent = textBeforeCursor;
      this.updateContent(contentId, textBeforeCursor);

      // Create new content
      const newContent: Content = {
        id: '', // Will be set by the backend
        noteId: this.note.id,
        data: textAfterCursor,
        type: 'text',
        version: 0,
        position: currentIndex + 1,
      };

      this.noteService.addContent(this.note.id, newContent, this.note.version).subscribe({
        next: (addedContent) => {
          // Focus on the new element after Angular renders it
          setTimeout(() => {
            const newParagraph = document.querySelector(`[data-content-id="${addedContent.id}"]`) as HTMLElement;
            if (newParagraph) {
              newParagraph.focus();
              const newRange = document.createRange();
              const newSelection = window.getSelection();
              newRange.setStart(newParagraph.childNodes[0] || newParagraph, 0);
              newRange.collapse(true);
              newSelection?.removeAllRanges();
              newSelection?.addRange(newRange);
            }
          }, 0);
        },
        error: (err) => {
          console.error('Error adding content', err);
          // Revert the change on error
          paragraph.textContent = originalText;
        },
      });
    } else if (event.key === 'Backspace') {
      const selection = window.getSelection();
      if (!selection || !selection.rangeCount) return;

      const { startOffset } = selection.getRangeAt(0);

      if (startOffset === 0) {
        if (!contentId || !this.note) return;

        const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
        if (currentIndex > 0) {
          event.preventDefault();
          const previousContent = this.note.contents[currentIndex - 1];
          const currentContent = this.note.contents[currentIndex];
          const mergedData = previousContent.data + (paragraph.textContent || '');
          const cursorOffset = previousContent.data.length;
          console.log('Merged content:', mergedData);

          this.updateContent(previousContent.id, mergedData);
          this.noteService.deleteContent(this.note.id, currentContent.id, this.note.version, currentContent.version).subscribe({
            next: () => {
              if (this.note) {
                setTimeout(() => {
                  const prevParagraph = document.querySelector(`[data-content-id="${previousContent.id}"]`) as HTMLElement;
                  if (prevParagraph) {
                    prevParagraph.focus();
                    const newRange = document.createRange();
                    const newSelection = window.getSelection();
                    const textNode = prevParagraph.childNodes[0] || prevParagraph;
                    // Ensure cursor position is not out of bounds
                    const newOffset = Math.min(cursorOffset, (textNode.textContent || '').length);
                    newRange.setStart(textNode, newOffset);
                    newRange.collapse(true);
                    newSelection?.removeAllRanges();
                    newSelection?.addRange(newRange);
                  }
                }, 0);
              }
            },
            error: (err) => {
              console.error('Error deleting content', err);
            }
          });
        }
      }
    }
  }

  private updateContent(contentId: string, newText: string): void {
    console.log('updateContent called for contentId:', contentId, 'newText:', newText);
    if (!this.note) {
      console.log('updateContent: Note is null, returning.');
      return;
    }

    const contentIndex = this.note.contents.findIndex(c => c.id === contentId);
    if (contentIndex === -1) {
      console.log('updateContent: Content not found, returning.');
      return;
    }

    const originalContent = this.note.contents[contentIndex];
    if (originalContent.data === newText) {
      console.log('updateContent: Content data is the same, returning.');
      return;
    }
    console.log('Original content:', originalContent);
    const updatedContent: Content = { ...originalContent, data: newText };
    console.log('Calling noteService.updateContent with:', updatedContent, 'note id:', updatedContent.noteId);

    this.noteService.updateContent(updatedContent).subscribe({
      next: () => {
        console.log('noteService.updateContent: next callback triggered.');
        if (this.note) {
          // this.note.contents[contentIndex] = { ...updatedContent, version: originalContent.version + 1 };
          console.log(`Content ${contentId} updated successfully`);
        }
      },
      error: (err) => {
        console.error(`noteService.updateContent: Error updating content ${contentId}`, err);
        // Revert the UI change on error
        const paragraph = document.querySelector(`[data-content-id="${contentId}"]`);
        if (paragraph) {
          paragraph.textContent = originalContent.data;
        }
      },
    });
  }
}
