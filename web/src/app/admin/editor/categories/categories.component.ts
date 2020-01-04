import { Component, OnInit } from '@angular/core';
import { faSearch, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { BehaviorSubject } from 'rxjs';
import { debounceTime, map } from 'rxjs/operators';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

@Component({
  selector: 'app-post-categories-editor',
  templateUrl: './categories.component.html',
  styleUrls: ['./categories.component.scss'],
})
export class PostCategoriesEditorComponent extends AbstractPostEditorComponent implements OnInit {

  /**
   * List of categories or tags to-be rendered as sidebar menu-item(s)
   */
  archives: Array<Category | Tag>;

  /**
   * A keyword for searching on the name of the archive
   */
  keyword: string;

  /**
   * Use to deboucing keypress event on search value
   */
  keyword$: BehaviorSubject<String> = new BehaviorSubject('');

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faSearch,
  };

  ngOnInit(): void {
    this.getAll('categories');

    this.keyword$.pipe(debounceTime(1600)).subscribe((search: string): void => {
      this.findAllArchives(search);
    });
  }

  onChange(): void {
    this.keyword$.next(this.keyword);
  }

  onKeyPress(): void {
    this.keyword$.next(this.keyword);
  }

  protected getAll(type: string): void {
    this.apollo.query({
      query: gql`
        {
          archives: ${type} { name slug }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ archives: Array<Category | Tag> }>): Array<Category> => result.data.archives),
    ).subscribe((archives: Array<Category | Tag>): void => {
      this.archives = archives;
    });
  }

  private findAllArchives(keyword: string): void {

  }
}
