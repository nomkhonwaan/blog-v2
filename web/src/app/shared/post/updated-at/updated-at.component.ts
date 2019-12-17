import { Component, OnInit } from '@angular/core';
import { AbstractPostComponent } from '../abstract-post.component';

@Component({
  selector: 'app-post-updated-at',
  templateUrl: './updated-at.component.html',
  styleUrls: ['./updated-at.component.scss'],
})
export class PostUpdatedAtComponent extends AbstractPostComponent implements OnInit {

  /**
   * A date-time that the post was updated
   */
  updatedAt: string;

  ngOnInit(): void {
    this.updatedAt = this.isZero(new Date(this.post.updatedAt)) ? this.post.createdAt : this.post.updatedAt;
  }

  isZero(date: Date): boolean {
    return date.getFullYear() === 1;
  }

}
