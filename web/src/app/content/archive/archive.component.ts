import { ChangeDetectorRef, ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';

@Component({
  selector: 'app-archive',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './archive.component.html',
  styleUrls: ['./archive.component.scss'],
})
export class ArchiveComponent implements OnInit {

  /**
   * When this property has been set to "all",
   * will retrieve the latest published posts from all categories
   * rather a single category or tag.
   */
  @Input()
  from: string;

  /**
   * List of posts
   */
  posts: Array<Post>;

  /**
   * Next and previous page URL will be displayed based on conditions
   */
  nextPage: Array<string>;
  previousPage: Array<string>;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private title: Title,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    if (!this.from) {
      this.from = (this.route.snapshot.data as { from: string }).from;
    }
  }

}
