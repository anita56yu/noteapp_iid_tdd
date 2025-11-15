import { Component, Input, Output, EventEmitter, OnChanges, SimpleChanges, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NoteService, Note, Content } from '../../services/note-service';
import { WebSocketService } from '../../services/websocket-service';
import { Subscription } from 'rxjs';

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
        console.log(contentToUpdate);
        if (contentToUpdate && message.content_version > contentToUpdate.version) {
          contentToUpdate.data = message.data;
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

  onContentInput(event: Event): void {
    // No debouncing needed, updates are triggered by blur or enter key
  }

  //BUG: When clicking outside and inside again it adds text automatically to the end.
  onContentBlur(): void {
    const paragraph = window.getSelection()?.anchorNode?.parentElement;
    if (paragraph && paragraph.tagName === 'P') {
      const contentId = paragraph.getAttribute('data-content-id');
      setTimeout(() => {
        if (contentId) {
          this.updateContent(contentId, paragraph.innerText);
        }
      }, 0);
    }
  }

  onContentKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      event.preventDefault();

      const selection = window.getSelection();
      if (!selection || !selection.rangeCount) return;

      const range = selection.getRangeAt(0);
      const { startContainer, startOffset } = range;

      let currentParagraph: Node | null = startContainer;
      if (currentParagraph.nodeType === Node.TEXT_NODE) {
        currentParagraph = currentParagraph.parentElement;
      }

      if (!currentParagraph || (currentParagraph as HTMLElement).tagName !== 'P' || !this.note) return;
      
      const contentId = (currentParagraph as HTMLElement).getAttribute('data-content-id');
      if (!contentId) return;

      const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
      if (currentIndex === -1) return;

      const originalText = (currentParagraph as HTMLElement).textContent || '';
      const textBeforeCursor = originalText.substring(0, startOffset);
      const textAfterCursor = originalText.substring(startOffset);

      // Update current content
      console.log('Updating content ID:', contentId, 'with text:', textBeforeCursor);
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
          // newContent.id = addedContent.id;
          // this.note?.contents.splice(currentIndex + 1, 0, newContent);
          // if (this.note) {
          //   this.note.version++;
          // }
          
          // Set focus on the new element after Angular renders it
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
          if (this.note) {
            this.note.contents[currentIndex].data = originalText;
          }
        },
      });
    } else if (event.key === 'Backspace') {
      const selection = window.getSelection();
      if (!selection || !selection.rangeCount) return;

      const range = selection.getRangeAt(0);
      const { startContainer, startOffset } = range;

      if (startOffset === 0) {
        let currentParagraph: Node | null = startContainer;
        if (currentParagraph.nodeType === Node.TEXT_NODE) {
          currentParagraph = currentParagraph.parentElement;
        }

        if (!currentParagraph || (currentParagraph as HTMLElement).tagName !== 'P' || !this.note) return;

        const contentId = (currentParagraph as HTMLElement).getAttribute('data-content-id');
        if (!contentId) return;

        const currentIndex = this.note.contents.findIndex(c => c.id === contentId);
        if (currentIndex > 0) {
          event.preventDefault();
          const previousContent = this.note.contents[currentIndex - 1];
          const currentContent = this.note.contents[currentIndex];
          const mergedData = previousContent.data + currentContent.data;
          var cursorOffset = previousContent.data.length;

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
                    newRange.setStart(textNode, cursorOffset);
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
