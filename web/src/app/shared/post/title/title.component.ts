import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { AbstractPostComponent } from '../abstract-post.component';

@Component({
  selector: 'app-post-title',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './title.component.html',
  styleUrls: ['./title.component.scss'],
})
export class PostTitleComponent extends AbstractPostComponent implements OnInit {

  /**
   * Use to indiciate whether link should enable or not
   */
  @Input()
  link = true;

  /**
   * When this proprety true; a generated link will go to editor page rather than single page
   */
  @Input()
  goToEditor = false;

  /**
   * An array of string to be composed to router link
   */
  href: string[];

  ngOnInit(): void {
    let date: Date;

    if (this.post.status === 'PUBLISHED') {
      date = this.toLocaleDate(this.post.publishedAt);
    } else {
      date = this.toLocaleDate(this.post.createdAt);
    }

    if (this.link) {
      this.href = [
        '/',
        date.getFullYear().toString(),
        (date.getMonth() + 1).toString(), // a month number start at 0 not 1
        date.getDate().toString(),
        this.post.slug,
        this.goToEditor ? 'edit' : undefined,
      ].filter((val: string): string => val);
    }
  }

  isDisabledLink(): boolean {
    return !!!this.link;
  }

  toLocaleDate(date: string): Date {
    return new Date(new Date(date).toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }));
  }

}
